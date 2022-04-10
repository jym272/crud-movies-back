package main

import (
	"backend/models"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

func (app *Application) getOneMovie(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Get the movie ID from the route
	//movieID := r.URL.Query().Get("id")
	movieID := ps.ByName("id")
	//convert string to int64
	if id, err := strconv.ParseInt(movieID, 10, 64); err == nil {
		app.logger.Println("getOneMovie: " + movieID)
		//create a dummy movie
		movie := models.Movie{
			ID:          id,
			Title:       "The Godfather",
			Year:        1972,
			ReleaseDate: time.Date(1972, time.January, 24, 0, 0, 0, 0, time.UTC),
			Runtime:     175,
			MPAARating:  "R",
			Rating:      9.2,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Director:    "Francis Ford Coppola",
			Poster:      "https://m.media-amazon.com/images/M/MV5BM2MyNjYxNmUtYTAwNi00MTYxLWJmNWYtYzZlODY3ZTk3OTFlXkEyXkFqcGdeQXVyNzkwMjQ5NzM@._V1_SX300.jpg",
		}
		err = app.writeJSON(w, http.StatusOK, movie, "movie")
		if err != nil {
			app.logger.Println("getOneMovie: " + err.Error())
			//app.writeError(w, http.StatusInternalServerError, err)
		}
	} else {
		app.logger.Println("getOneMovie: " + err.Error())
	}

}

//get all movies
func (app *Application) getAllMovies(w http.ResponseWriter, r *http.Request) {
	//movies, err := app.movieStore.GetAll()
	//if err != nil {
	//	app.clientError(w, http.StatusNotFound)
	//	return
	//}
	//
	//// Render the template
	//app.render(w, r, "movies.page.tmpl", &templateData{
	//	Movies: movies,
	//})
}
