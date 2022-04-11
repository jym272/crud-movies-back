package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (app *Application) getOneMovie(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the movie ID from the route
	//movieID := r.URL.Query().Get("id")
	movieID := ps.ByName("id")
	//convert string to int64
	if id, err := strconv.ParseInt(movieID, 10, 64); err == nil {
		movie, err := app.models.DB.GetMovie(int(id))
		if err != nil {
			app.errorJSON(w, http.StatusNotFound, err)
			app.logger.Println("getOneMovie: " + err.Error())
			//http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		err = app.writeJSON(w, http.StatusOK, movie, "movie")
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("getOneMovie: " + err.Error())
		}
	} else {
		app.errorJSON(w, http.StatusBadRequest, err)
		app.logger.Println("getOneMovie: " + err.Error())
	}

}

//get all movies
func (app *Application) getAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.models.DB.GetAll()
	if err != nil {
		app.errorJSON(w, http.StatusNotFound, err)
		app.logger.Println("getAllMovies: " + err.Error())
		//http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("getAllMovies: " + err.Error())
	}
}

//getGenres
func (app *Application) getGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.models.DB.GetGenres()
	if err != nil {
		app.errorJSON(w, http.StatusNotFound, err)
		app.logger.Println("getGenres: " + err.Error())
		//http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	err = app.writeJSON(w, http.StatusOK, genres, "genres")
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("getGenres: " + err.Error())
	}
}
