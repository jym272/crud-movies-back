package main

import (
	"backend/models"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

var validUser = models.User{
	ID:       10,
	Name:     "John Doe",
	Email:    "jym272@gmail.com",
	Password: generateHasPassword("password"),
}

func generateHasPassword(password string) string {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "password"
	}
	return string(fromPassword)

}

func (app *Application) signinHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, err)
		app.logger.Println("signinHandler1: " + err.Error())
		return
	}
	//TODO: later we will check the user in the database

	if credentials.Username != validUser.Email {
		app.errorJSON(w, http.StatusUnauthorized, err) //invalid credentials
		app.logger.Println("signinHandler2: " + err.Error())
		return
	}

	hashedPassword := validUser.Password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
	if err != nil {
		app.errorJSON(w, http.StatusUnauthorized, err) //invalid credentials
		app.logger.Println("signinHandler3: " + err.Error())
		return
	}

	token, err := app.createToken(validUser)
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

func (app *Application) createToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["name"] = user.Name
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["iat"] = time.Now().Unix()
	claims["nbf"] = time.Now().Unix()
	claims["iss"] = "mydomain.com" //http://localhost:8080
	claims["aud"] = []string{"mydomain.com"}

	return token.SignedString([]byte(app.config.secretKey))
}
