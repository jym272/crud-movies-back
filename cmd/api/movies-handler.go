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
	adjacent := r.URL.Query().Get("adjacent_ids")
	type adjacentIds struct {
		Ids [2]int64 `json:"ids"`
	}

	type ResponseType struct {
		Movie         *models.Movie `json:"movie"`
		AdjacentIds   adjacentIds   `json:"adjacent_ids"`
		WithGenreName string        `json:"with_genre_name"`
	}

	movieID := ps.ByName("id")
	//convert string to int64
	if id, err := strconv.ParseInt(movieID, 10, 64); err == nil {
		var response ResponseType

		if adjacent == "true" {
			//get adjacents ids in the db

			withGenre := r.URL.Query().Get("withgenre")
			var genreID int64
			if withGenre != "" {
				genreID, err = strconv.ParseInt(withGenre, 10, 64)
				if err != nil {
					app.errorJSON(w, http.StatusInternalServerError, err)
					app.logger.Println("getOneMovie: " + err.Error())
					return
				}
				//get genre name with id
				genreNameByID, err := app.models.DB.GetGenreNameByID(genreID)
				if err != nil {
					app.errorJSON(w, http.StatusInternalServerError, err)
					app.logger.Println("getOneMovie: " + err.Error())
					return
				}
				response.WithGenreName = genreNameByID
			}

			var adjacent adjacentIds
			ids, err := app.models.DB.GetMoviesIds(genreID)
			if err != nil {
				app.errorJSON(w, http.StatusNotFound, err)
				app.logger.Println("getOneMovie0: " + err.Error())
				return
			}
			//find the adjacent ids of id in the ids array
			for i := 0; i < len(ids); i++ {
				if ids[i] == id {
					if i-1 >= 0 {
						adjacent.Ids[0] = ids[i-1]
					} else {
						//the adjacent is the last movie
						adjacent.Ids[0] = ids[len(ids)-1]
					}
					if i+1 < len(ids) {
						adjacent.Ids[1] = ids[i+1]
					} else {
						//the adjacent is the first movie
						adjacent.Ids[1] = ids[0]
					}
					break
				}
			}
			response.AdjacentIds = adjacent
		}

		movie, err := app.models.DB.GetMovie(id)
		if err != nil {
			app.errorJSON(w, http.StatusNotFound, err)
			app.logger.Println("getOneMovie: " + err.Error())
			//http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		response.Movie = movie
		err = app.writeJSON(w, http.StatusOK, response, "")
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
		//get adjacent genres ids
		adjacentGenresIdsQuery := r.URL.Query().Get("adjacent_genres_ids")
		type adjacent struct {
			Id   int64  `json:"id"`
			Name string `json:"name"`
		}
		type adjacentGenresType struct {
			Next *adjacent `json:"next"`
			Prev *adjacent `json:"previous"`
		}
		//make ajacentGenresType
		var adjacentGenres = &adjacentGenresType{}

		if adjacentGenresIdsQuery == "true" {
			genres, err = app.models.DB.GetGenres()
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				app.logger.Println("getAllMovies: " + err.Error())
				return
			}
			//find the adjacent ids of id in the ids array
			for i := 0; i < len(genres); i++ {
				if genres[i].ID == parseInt {
					if i-1 >= 0 {
						adjacentGenres.Prev = &adjacent{genres[i-1].ID, genres[i-1].Name}
					} else {
						//the adjacent is the last movie
						adjacentGenres.Prev = &adjacent{genres[len(genres)-1].ID, genres[len(genres)-1].Name}
					}
					if i+1 < len(genres) {
						adjacentGenres.Next = &adjacent{genres[i+1].ID, genres[i+1].Name}
					} else {
						//the adjacent is the first movie
						adjacentGenres.Next = &adjacent{genres[0].ID, genres[0].Name}
					}
					break
				}
			}
			err = app.writeJSON(w, http.StatusOK, &struct {
				GenreName      string              `json:"genre_name"`
				Movies         []*models.Movie     `json:"movies"`
				AdjacentGenres *adjacentGenresType `json:"adjacent_genres"`
			}{genreName, movies, adjacentGenres}, "")
		} else {
			err = app.writeJSON(w, http.StatusOK, &struct {
				GenreName string          `json:"genre_name"`
				Movies    []*models.Movie `json:"movies"`
			}{genreName, movies}, "")
		}

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

func (app *Application) favoritesHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userId").(int64)
	movieIDQuery := r.URL.Query().Get("movie")
	action := r.URL.Query().Get("action")

	if action == "retrievefavorites" {
		type MoviesIDs struct {
			MovieIdsArray []int64 `json:"ids"`
		}
		var movieIDs MoviesIDs
		err := json.NewDecoder(r.Body).Decode(&movieIDs)
		if err != nil {
			app.errorJSON(w, http.StatusBadRequest, err)
			app.logger.Println("favoritesHandler1: " + err.Error())
			return
		}
		result := make(map[int64]bool, len(movieIDs.MovieIdsArray))
		for _, movieID := range movieIDs.MovieIdsArray {
			isFav := app.models.DB.IsFav(movieID, userID)
			result[movieID] = isFav

		}
		err = app.writeJSON(w, http.StatusOK, result, "")
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("favoritesHandler2: " + err.Error())
		}
		return
	}

	if action == "list" {
		//get all the favorites of the user
		movies_, err := app.models.DB.GetFavorites(userID)
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("favoritesHandler1: " + err.Error())
			return
		}
		err = app.writeJSON(w, http.StatusOK, movies_, "movies")
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("getAllMovies: " + err.Error())
		}
		return
	}

	parseMovieID, err := strconv.ParseInt(movieIDQuery, 10, 64)
	if err != nil {
		app.errorJSON(w, http.StatusBadRequest, err)
		app.logger.Println("addToFav1: " + err.Error())
		return
	}
	isFav := app.models.DB.IsFav(parseMovieID, userID)
	//switch action
	switch action {
	case "query":
		var data string
		if isFav {
			data = "favorite"
		} else {
			data = "not favorite"
		}
		err = app.writeJSON(w, http.StatusOK, data, "")
	case "add":
		//add to favorite_movies table
		if !isFav {
			err = app.models.DB.AddToFav(parseMovieID, userID)
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				app.logger.Println("addToFav3: " + err.Error())
				return
			}
		}
		err = app.writeJSON(w, http.StatusOK, "added", "")
	case "remove":
		//remove from favorite_movies table
		if isFav {
			err = app.models.DB.RemoveFromFav(parseMovieID, userID)
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				app.logger.Println("addToFav5: " + err.Error())
				return
			}
		}
		err = app.writeJSON(w, http.StatusOK, "removed", "")
	default:
		err = app.writeJSON(w, http.StatusBadRequest, "invalid action", "error")
	}
	if err != nil { //repeat code
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("addToFav7: " + err.Error())
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
