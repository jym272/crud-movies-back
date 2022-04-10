package models

import "time"

type minutes int

type Movie struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Year        int          `json:"year"`
	Rating      float32      `json:"rating"`
	ReleaseDate time.Time    `json:"release_date"`
	Runtime     minutes      `json:"runtime"`
	MPAARating  string       `json:"mpaa_rating"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	MovieGenres []MovieGenre `json:"-"`
	Poster      string       `json:"poster"`
	Director    string       `json:"director"`
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
