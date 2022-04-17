package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func (app *Application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			//anonymous user can have protected routes
			app.errorJSON(w, http.StatusUnauthorized, errors.New("missing auth header"))
			return
		}
		headerParts := strings.Split(authHeader, " ")

		if len(headerParts) != 2 {
			app.errorJSON(w, http.StatusBadRequest, errors.New("invalid auth header"))
			return
		}
		bearer := headerParts[0]
		token := headerParts[1]

		if bearer != "Bearer" {
			app.errorJSON(w, http.StatusBadRequest, errors.New("invalid auth header - bearer not found"))
			return
		}

		token_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(app.config.secretKey), nil
		})
		if err != nil {
			app.errorJSON(w, http.StatusUnauthorized, errors.New("unauthorized - failed to parse token"))
			return
		}

		if token_.Valid {
			if claims, ok := token_.Claims.(jwt.MapClaims); ok {
				fmt.Println(claims["id"], claims["nbf"])
				//TODO: se peude usar extraer informacion del token respecto al usuario y router acorde a eso
			} else {
				fmt.Println("Error parsing claims")
			}

			next.ServeHTTP(w, r)
		} else if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				app.errorJSON(w, http.StatusUnauthorized, errors.New("unauthorized - malformed token"))
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				app.errorJSON(w, http.StatusUnauthorized, errors.New("unauthorized - token expired"))
				return
			} else if ve.Errors&(jwt.ValidationErrorAudience) != 0 {
				app.errorJSON(w, http.StatusUnauthorized, errors.New("unauthorized - invalid audience"))
				return
			} else if ve.Errors&(jwt.ValidationErrorIssuer) != 0 {
				app.errorJSON(w, http.StatusUnauthorized, errors.New("unauthorized - invalid issuer"))
				return
			} else {
				app.errorJSON(w, http.StatusUnauthorized, errors.New("unauthorized - failed to parse token"))
				return
			}
		} else {
			app.errorJSON(w, http.StatusUnauthorized, errors.New("unauthorized - failed to parse token"))
			return
		}

	})
}
