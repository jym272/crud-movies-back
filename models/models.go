package models

import (
	"database/sql"
	"time"
)

//Models is the wrapper for database
type Models struct {
	DB DBModel
}

// NewModels returns a new instance of Models, db pool
func NewModels(db *sql.DB) *Models {
	return &Models{
		DB: DBModel{
			DB: db,
		},
	}
}

//Movie es the type for the movie
type Movie struct {
	ID                int64                  `json:"id"`
	Title             string                 `json:"title"`
	Description       string                 `json:"description"`
	Year              int                    `json:"year"`
	ReleaseDate       time.Time              `json:"release_date"`
	Runtime           int                    `json:"runtime"`
	Rating            float32                `json:"rating"`
	MPAARating        string                 `json:"mpaa_rating"`
	CreatedAt         time.Time              `json:"-"`
	UpdatedAt         time.Time              `json:"-"`
	MovieGenres       map[int]string         `json:"genres"`
	Genres            []*Genre               `json:"genres_list"`         //for graphql
	AdjacentMoviesIds *AdjacentMoviesIdsType `json:"adjacent_movies_ids"` //for graphql
	Poster            string                 `json:"poster"`
	UserID            int64                  `json:"user_id"`
}

//&models.AdjacentMoviesIds{Previous: Ids[0], Next: Ids[1]}

type AdjacentMoviesIdsType struct {
	Previous int64 `json:"previous"`
	Next     int64 `json:"next"`
}
type Genre struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type MovieGenre struct {
	ID        int64     `json:"-"`
	MovieID   int64     `json:"-"`
	GenreID   int64     `json:"-"`
	Genre     Genre     `json:"genre"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
