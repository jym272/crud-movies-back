package main

import (
	"backend/models"
	"errors"
	"github.com/graphql-go/graphql"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var movies []*models.Movie

var genreType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Genre",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

//graphql schema definition -> exposed to the client
var movieType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Movie",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"year": &graphql.Field{
				Type: graphql.Int,
			},
			"runtime": &graphql.Field{
				Type: graphql.Int,
			},
			"release_date": &graphql.Field{
				Type: graphql.DateTime,
			},
			"rating": &graphql.Field{
				Type: graphql.Int,
			},
			"mpaa_rating": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"genres_list": &graphql.Field{
				Type: graphql.NewList(genreType),
			},
		},
	},
)

var fields = graphql.Fields{
	"hello": &graphql.Field{
		Type: graphql.String,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return "world", nil
		},
	},
	"movie": &graphql.Field{
		Type:        movieType,
		Description: "Get movie by id",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, ok := p.Args["id"].(int)
			if ok {
				for _, movie := range movies {
					if movie.ID == int64(id) {
						genres := make([]*models.Genre, len(movie.MovieGenres))
						i := 0
						for id, genre := range movie.MovieGenres {
							genres[i] = &models.Genre{ID: int64(id), Name: genre}
							i++
						}
						movie.Genres = genres
						return movie, nil
					}
				}
			}
			return nil, nil
		},
	},
	"list": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "Get list of movies",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return movies, nil
		},
	},
	"search": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "Search movies by title",
		Args: graphql.FieldConfigArgument{
			"titleContains": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			search, ok := p.Args["titleContains"].(string)
			search = strings.TrimSpace(strings.ToLower(search))
			var results []*models.Movie
			if ok {
				for _, movie := range movies {
					if strings.Contains(strings.ToLower(movie.Title), search) {
						results = append(results, movie)
					}
				}
				return results, nil
			}
			return nil, nil
		},
	},
	//get movies by genre id
	"genre": &graphql.Field{
		Type:        graphql.NewList(movieType),
		Description: "Get movies by genre id",
		Args: graphql.FieldConfigArgument{
			"genreId": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			genreId, ok := p.Args["genreId"].(int)
			if ok {
				var results []*models.Movie
				//
				for _, movie := range movies {
					for id := range movie.MovieGenres {
						if genreId == id {
							results = append(results, movie)
						}
					}
				}
				return results, nil
			}
			return nil, nil
		},
	},
}

func (app *Application) moviesGraphQL(w http.ResponseWriter, r *http.Request) {
	var err error
	movies, err = app.models.DB.GetAll()
	if err != nil {
		app.errorJSON(w, http.StatusNotFound, err)
		app.logger.Println("moviesGraphQL1: " + err.Error())
		return
	}

	q, err := ioutil.ReadAll(r.Body) //manage err later
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("moviesGraphQL2: " + err.Error())
		return
	}
	query := string(q) //{genre(genreId: 3){id title}}

	var genreId int

	if strings.Contains(query, "genre(genreId:") {
		re := regexp.MustCompile(`genreId: (\d+)`)
		match := re.FindStringSubmatch(query)
		if len(match) > 1 {
			genreId, err = strconv.Atoi(match[1])
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				app.logger.Println("moviesGraphQL3: " + err.Error())
				return
			}
		}

	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("moviesGraphQL2: " + err.Error())
		return
	}

	params := graphql.Params{Schema: schema, RequestString: query}
	rJSON := graphql.Do(params)

	if len(rJSON.Errors) > 0 {
		errArray := make([]string, len(rJSON.Errors))
		for i, err := range rJSON.Errors {
			errArray[i] = err.Error()
		}
		listString := strings.Join(errArray, ",")
		app.errorJSON(w, http.StatusBadRequest, errors.New(listString))
		app.logger.Println("moviesGraphQL3: " + listString)
		return
	}
	if genreId > 0 { //asumo que una genre id nunca puede ser cero
		response := rJSON.Data.(map[string]interface{})
		moviesResultsArray := response["genre"].([]interface{})
		moviesResults := make(map[string][]interface{}, len(moviesResultsArray))

		var genreId int

		re := regexp.MustCompile(`genreId: (\d+)`)
		match := re.FindStringSubmatch(query)
		if len(match) > 1 {
			genreId, err = strconv.Atoi(match[1])
			if err != nil {
				app.errorJSON(w, http.StatusInternalServerError, err)
				app.logger.Println("moviesGraphQL3: " + err.Error())
				return
			}
		}
		//getGenreNameByID
		genreName, err := app.models.DB.GetGenreNameByID(int64(genreId))
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("moviesGraphQL3.1: " + err.Error())
			return
		}

		moviesResults[genreName] = moviesResultsArray

		err = app.writeJSON(w, http.StatusOK, moviesResults, "data")
		if err != nil {
			app.errorJSON(w, http.StatusInternalServerError, err)
			app.logger.Println("moviesGraphQL4: " + err.Error())
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, rJSON, "")
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err)
		app.logger.Println("moviesGraphQL4: " + err.Error())
	}
}
