package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"github.com/dlclark/regexp2"
	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

var validUser = models.User{
	ID:       10,
	Username: "jym272@gmail.com",
	Password: generateHashPassword("password"),
}

func generateHashPassword(password string) string {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
		return "password"
	}
	return string(fromPassword)

}

func cleanAndValidateCredentials(credentials *Credentials, app *Application, w *http.ResponseWriter) error {
	//Valid and clean data
	// Clean received data
	credentials.Username = strings.TrimSpace(strings.ToLower(credentials.Username))
	credentials.Password = strings.TrimSpace(credentials.Password)

	// http://emailregex.com/
	var regexEmail = `^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`
	re := regexp.MustCompile(regexEmail)
	if !re.MatchString(credentials.Username) {
		errorMsg := "invalid username"
		app.errorJSON(*w, http.StatusBadRequest, errors.New(errorMsg))
		return errors.New(errorMsg)
	}

	/*
	 * Check if newPassword is valid: at least 8 characters, at least one number, at least one lowercase and one uppercase letter.
	 * Special characters are allowed but not required, even spaces.
	 * https://stackoverflow.com/a/49721224/3681450
	 */
	var regexPassword = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$`
	rePass := regexp2.MustCompile(regexPassword, 0)
	isMatch, err := rePass.MatchString(credentials.Password)
	if err != nil {
		app.errorJSON(*w, http.StatusInternalServerError, err)
		return err
	}
	if isMatch == false {
		errorMsg := "password: at least 8 characters, at least one number, at least one lowercase and one uppercase letter. Special characters are allowed but not required, even spaces"
		app.errorJSON(*w, http.StatusBadRequest, errors.New(errorMsg))
		return errors.New(errorMsg)
	}
	return nil

}

func (app *Application) signinHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	withGoogle := r.URL.Query().Get("google")
	var user *models.User
	if withGoogle == "true" {
		//only username is required
		type GoogleCredentials struct {
			Username string `json:"email"`
		}
		var googleCredentials GoogleCredentials
		err := json.NewDecoder(r.Body).Decode(&googleCredentials)
		if err != nil {
			app.errorJSON(w, http.StatusBadRequest, err)
			return
		}
		//there is no need to clean and validate data
		//find the user in the database, if not found, create a new user
		user, err = app.models.DB.GetUser(googleCredentials.Username)
		if err != nil {
			//user not found, create a new user
			user = &models.User{
				Username: googleCredentials.Username,
				Password: "google",
			}
			err = app.models.DB.CreateUser(user)
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				return
			}
			//get user from db
			user, err = app.models.DB.GetUser(googleCredentials.Username)
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				app.logger.Println("signupHandler6: " + err.Error())
				return
			}
			//create a new token and send it to the client

		} // else: user found, create a new token and send it to the client

	} else { //credentials sign in
		var credentials Credentials
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			app.errorJSON(w, http.StatusBadRequest, err)
			app.logger.Println("signinHandler1: " + err.Error())
			return
		}
		err = cleanAndValidateCredentials(&credentials, app, &w)
		if err != nil {
			app.logger.Println("signinHandler2: " + err.Error())
			return
		}
		user, err = app.models.DB.GetUser(credentials.Username)
		if err != nil {
			errorMsg := "user not found"
			app.errorJSON(w, http.StatusUnauthorized, errors.New(errorMsg)) //invalid credentials
			app.logger.Println("signinHandler2: " + errorMsg)
			return
		}

		if user.Password == "google" {
			errorMsg := "google account, login with google or create a password with the signup form"
			app.errorJSON(w, http.StatusUnauthorized, errors.New(errorMsg)) //invalid credentials
			app.logger.Println("signinHandler2: " + errorMsg)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			errorMsg := "invalid password"
			app.errorJSON(w, http.StatusUnauthorized, errors.New(errorMsg)) //invalid credentials
			app.logger.Println("signinHandler3: " + err.Error())
			return
		}

	}

	token, err := app.createToken(user)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("signinHandler4: " + err.Error())
		return
	}
	err = app.writeJSON(w, http.StatusOK, token, "jwt")
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("signinHandler5: " + err.Error())
	}
}

func (app *Application) createToken(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["email"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = time.Now().Unix()
	claims["iss"] = "mydomain.com" //http://localhost:8080
	claims["aud"] = []string{"mydomain.com"}

	return token.SignedString([]byte(app.config.secretKey))
}

func (app *Application) signupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, err)
		app.logger.Println("signupHandler1: " + err.Error())
		return
	}
	err = cleanAndValidateCredentials(&credentials, app, &w)
	if err != nil {
		app.logger.Println("signupHandler2: " + err.Error())
		return
	}

	// Check if user already exists
	//un user es google o credentials, si es google, no tiene password, password es "google"
	//un nuevo user válido, no existe en el db y si existe su password es "google"

	var user *models.User
	user, err = app.models.DB.GetUser(credentials.Username)
	if err == nil { //user already exists
		if user.Password != "google" {
			errorMsg := "user already exists"
			app.errorJSON(w, http.StatusBadRequest, errors.New(errorMsg))
			app.logger.Println("signupHandler4: " + errorMsg)
			return
		} else {
			//user already exists, but is google, so update password
			err = app.models.DB.UpdateUserPasswordByUsername(credentials.Username, generateHashPassword(credentials.Password))
			if err != nil {
				errorMsg := "error updating user"
				app.errorJSON(w, http.StatusInternalServerError, errors.New(errorMsg))
				app.logger.Println("signupHandler5: " + errorMsg)
				return
			}
		}
	} else { //user does not exist.
		// Create new user
		user = &models.User{
			Username: credentials.Username,
			Password: generateHashPassword(credentials.Password),
		}

		err = app.models.DB.CreateUser(user)
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("signupHandler5: " + err.Error())
			return
		}
	}
	//get user from db
	user, err = app.models.DB.GetUser(credentials.Username)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("signupHandler6: " + err.Error())
		return
	}
	//send user to client
	err = app.writeJSON(w, http.StatusOK, user, "")
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("signinHandler7: " + err.Error())
	}
}
