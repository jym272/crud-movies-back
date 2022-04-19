package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"golang.org/x/net/context"
	"net/http"
)

func (app *Application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *Application) routes() http.Handler {
	router := httprouter.New()

	secure := alice.New(app.checkToken)

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.POST("/v1/signin", app.signinHandler)
	router.HandlerFunc(http.MethodPost, "/v1/graphql", app.moviesGraphQL)
	router.GET("/v1/admin/delete", app.wrap(secure.ThenFunc(app.deleteOneMovie)))

	router.PUT("/v1/admin/movie", app.wrap(secure.ThenFunc(app.editOneMovie))) //update or create a movie

	router.GET("/v1/movie/:id", app.getOneMovie)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getGenres)

	return app.enableCORS(router)
}
