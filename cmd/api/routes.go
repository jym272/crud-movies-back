package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *Application) routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	//router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.GET("/v1/movie/:id", app.getOneMovie)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)

	//router.GET("/api/v1/users", app.getUsers)
	//router.GET("/api/v1/users/:id", app.getUser)
	//router.POST("/api/v1/users", app.createUser)
	//router.PUT("/api/v1/users/:id", app.updateUser)
	//router.DELETE("/api/v1/users/:id", app.deleteUser)
	return router
}
