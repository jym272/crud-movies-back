package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.GET("/v1/admin/delete", app.deleteOneMovie)
	router.PUT("/v1/admin/movie", app.editOneMovie) //update or create a movie
	router.GET("/v1/movie/:id", app.getOneMovie)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getGenres)

	return app.enableCORS(router)
}
