package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	//router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)
	router.GET("/v1/movie/:id", app.getOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getGenres)
	//router.GET("/v1/movies?genre=", app.getMoviesByGenre) //returns all movies with a specific genre

	//router.GET("/api/v1/users", app.getUsers)
	//router.GET("/api/v1/users/:id", app.getUser)
	//router.POST("/api/v1/users", app.createUser)
	//router.PUT("/api/v1/users/:id", app.updateUser)
	//router.DELETE("/api/v1/users/:id", app.deleteUser)
	return app.enableCORS(router)
}
