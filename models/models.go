package models

import (
	"database/sql"
	"time"
)

type minutes int

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
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Year        int          `json:"year"`
	ReleaseDate time.Time    `json:"release_date"`
	Runtime     minutes      `json:"runtime"`
	Rating      float32      `json:"rating"`
	MPAARating  string       `json:"mpaa_rating"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	MovieGenres []MovieGenre `json:"-"`
}

type Genre struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MovieGenre struct {
	ID        int64     `json:"id"`
	MovieID   int64     `json:"movie_id"`
	GenreID   int64     `json:"genre_id"`
	Genre     Genre     `json:"genre"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
