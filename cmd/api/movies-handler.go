package main

import (
	"backend/models"
	"encoding/json"
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
		movie, err := app.models.DB.GetMovie(id)
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
	genreIDQuery := r.URL.Query().Get("genre_id")
	if genreIDQuery != "" {
		parseInt, err := strconv.ParseInt(genreIDQuery, 10, 64)
		if err != nil {
			app.errorJSON(w, http.StatusBadRequest, err)
			app.logger.Println("getAllMovies: " + err.Error())
			return
		}
		movies, err := app.models.DB.GetMoviesByGenreWithID(parseInt)
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("getAllMovies: " + err.Error())
			return
		}
		//getGenreNameByID
		genreName, err := app.models.DB.GetGenreNameByID(parseInt)
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("getAllMovies: " + err.Error())
			return
		}

		err = app.writeJSON(w, http.StatusOK, movies, genreName)
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("getAllMovies: " + err.Error())
		}

	} else {

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

func (app *Application) editOneMovie(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	type response struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}
	//extract payload from the request
	var movie models.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)

	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, err)
		app.logger.Println("editOneMovie: " + err.Error())
		return
	}
	//extract year from movie.release_date
	year, _, _ := movie.ReleaseDate.Date()
	movie.Year = year

	movieIDQuery := r.URL.Query().Get("id")
	if movieIDQuery != "" {
		parseMovieID, err := strconv.ParseInt(movieIDQuery, 10, 64)
		if err != nil {
			app.errorJSON(w, http.StatusBadRequest, err)
			app.logger.Println("editOneMovie: " + err.Error())
			return
		}
		//check if the movie exists
		exist, err := app.models.DB.MovieExists(parseMovieID)
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("editOneMovie: " + err.Error())
			return
		}
		if exist {
			//update the movie
			err = app.models.DB.UpdateMovie(parseMovieID, &movie)
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				app.logger.Println("editOneMovie: " + err.Error())
				return
			}

			_response := response{
				Message: "Movie updated",
				Status:  http.StatusOK,
			}
			err = app.writeJSON(w, http.StatusOK, _response, "")
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				app.logger.Println("getGenres: " + err.Error())
			}
			return
		}

	} //create new movie: if there is no movieID in the query or if the movieID is not found in the database, or if the movieID 0 ->this id does not exist
	err = app.models.DB.InsertMovie(&movie)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("editOneMovie: " + err.Error())
		return
	}
	_response := response{
		Message: "Movie created",
		Status:  http.StatusOK,
	}
	err = app.writeJSON(w, http.StatusOK, _response, "")
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("editOneMovie: " + err.Error())
	}
}
