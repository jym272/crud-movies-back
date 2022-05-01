package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

func (app *Application) editOneMovie(w http.ResponseWriter, r *http.Request) {

	// read userID from the request context
	userID := r.Context().Value("userId").(int64)

	//TODO:later compera the owner of the movie with the userID

	type response struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}
	//extract payload from the request
	var movie models.Movie
	err := json.NewDecoder(r.Body).Decode(&movie)

	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, err)
		app.logger.Println("editOneMovie1: " + err.Error())
		return
	}

	//extract year from movie.release_date
	year, _, _ := movie.ReleaseDate.Date()
	movie.Year = year

	//add a poster to the movie
	movie.Poster = getPoster(movie.Title)

	movieIDQuery := r.URL.Query().Get("id")
	if movieIDQuery != "" {
		parseMovieID, err := strconv.ParseInt(movieIDQuery, 10, 64)
		if err != nil {
			app.errorJSON(w, http.StatusBadRequest, err)
			app.logger.Println("editOneMovie2: " + err.Error())
			return
		}
		//check if the movie exists
		exist, err := app.models.DB.MovieExists(parseMovieID)
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("editOneMovie3: " + err.Error())
			return
		}
		if exist {
			//check if the user is the owner of the movie
			isOwner, err := app.models.DB.IsOwner(parseMovieID, userID)
			if !isOwner || err != nil {
				app.errorJSON(w, http.StatusUnauthorized, errors.New("you are not the owner of this movie"))
				app.logger.Println("editOneMovie4: " + err.Error())
				return
			}

			//update the movie
			err = app.models.DB.UpdateMovie(parseMovieID, &movie)
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				app.logger.Println("editOneMovie5: " + err.Error())
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
	////prevent movie title duplication
	movieExist, err := app.models.DB.MovieTitleExist(movie.Title)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("editOneMovie: " + err.Error())
		return
	}
	if movieExist {
		msg := "movie title already exist"
		app.errorJSON(w, http.StatusBadRequest, errors.New(msg))
		return
	}

	//update the userId to the movie
	movie.UserID = userID
	err = app.models.DB.InsertMovie(&movie)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("editOneMovie6: " + err.Error())
		return
	}
	_response := response{
		Message: "Movie created",
		Status:  http.StatusOK,
	}
	err = app.writeJSON(w, http.StatusOK, _response, "")
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("editOneMovie7: " + err.Error())
	}
}

func (app *Application) getMyMovies(w http.ResponseWriter, r *http.Request) {
	// read userID from the request context
	userID := r.Context().Value("userId").(int64)

	//get the movies of the user
	movies, err := app.models.DB.GetAll(userID)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("getMyMovies1: " + err.Error())
		return
	}
	err = app.writeJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("getMyMovies2: " + err.Error())
	}
}

func (app *Application) deleteOneMovie(w http.ResponseWriter, r *http.Request) {

	// read userID from the request context
	userID := r.Context().Value("userId").(int64)

	type response struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}
	movieIDQuery := r.URL.Query().Get("id")
	if movieIDQuery != "" {
		parseMovieID, err := strconv.ParseInt(movieIDQuery, 10, 64)
		if err != nil {
			app.errorJSON(w, http.StatusBadRequest, err)
			app.logger.Println("deleteOneMovie1: " + err.Error())
			return
		}
		//check if the movie exists
		exist, err := app.models.DB.MovieExists(parseMovieID)
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("deleteOneMovie2: " + err.Error())
			return
		}
		if exist {
			//check if the user is the owner of the movie
			isOwner, err := app.models.DB.IsOwner(parseMovieID, userID)
			if !isOwner || err != nil {
				app.errorJSON(w, http.StatusUnauthorized, errors.New("you are not the owner of this movie"))
				app.logger.Println("deleteOneMovie4: " + err.Error())
				return
			}
			//delete the movie
			err = app.models.DB.DeleteMovie(parseMovieID)
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				app.logger.Println("deleteOneMovie5: " + err.Error())
				return
			}
			_response := response{
				Message: "Movie deleted",
				Status:  http.StatusOK,
			}
			err = app.writeJSON(w, http.StatusOK, _response, "")
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				app.logger.Println("deleteOneMovie6: " + err.Error())
			}
			return
		}

	}
	_response := response{
		Message: "Movie not found",
		Status:  http.StatusNotFound,
	}
	err := app.writeJSON(w, http.StatusNotFound, _response, "")
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("deleteOneMovie7: " + err.Error())
	}

}

//refactorizando la funcion
func getPoster(movieTitle string) string {
	type MovieDBType struct {
		Page    int `json:"page"`
		Results []struct {
			Adult            bool    `json:"adult"`
			BackdropPath     *string `json:"backdrop_path"`
			GenreIds         []int   `json:"genre_ids"`
			Id               int     `json:"id"`
			OriginalLanguage string  `json:"original_language"`
			OriginalTitle    string  `json:"original_title"`
			Overview         string  `json:"overview"`
			Popularity       float64 `json:"popularity"`
			PosterPath       string  `json:"poster_path"`
			ReleaseDate      string  `json:"release_date"`
			Title            string  `json:"title"`
			Video            bool    `json:"video"`
			VoteAverage      float64 `json:"vote_average"`
			VoteCount        int     `json:"vote_count"`
		} `json:"results"`
		TotalPages   int `json:"total_pages"`
		TotalResults int `json:"total_results"`
	}

	client := &http.Client{}

	apiKey := os.Getenv("API_KEY_MOVIE_DB")

	//https://api.themoviedb.org/3/search/movie?api_key=f6646a0386887b9fd168de141c70bd9b&query=the%20shawshank%20redemption
	title := url.QueryEscape(movieTitle)
	urlString := "https://api.themoviedb.org/3/search/movie?api_key=" + apiKey + "&query=" + title
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		fmt.Println(err)
		return "" //return the movie without poster
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "" //return the movie without poster
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("getPoster", err)
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "" //return the movie without poster
	}
	var movieDB MovieDBType
	err = json.Unmarshal(body, &movieDB)
	if err != nil {
		fmt.Println(err)
		return "" //return the movie without poster
	}
	if len(movieDB.Results) > 0 {
		return movieDB.Results[0].PosterPath
	}
	return ""
}
