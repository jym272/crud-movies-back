package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func getGenres(m *DBModel, movieId int, ctx context.Context) (*map[int]string, error) {
	query := "SELECT mg.id, mg.movie_id, mg.genre_id, genres.genre_name FROM genres INNER JOIN movies_genres mg ON genres.id = mg.genre_id WHERE mg.movie_id = $1"
	rows, err := m.DB.QueryContext(ctx, query, movieId)
	if err != nil {
		return nil, err
	}
	movieGenres := make(map[int]string)
	for rows.Next() {
		var genre MovieGenre
		err = rows.Scan(&genre.ID, &genre.MovieID, &genre.GenreID, &genre.Genre.Name)
		if err != nil {
			return nil, err
		}
		movieGenres[int(genre.GenreID)] = genre.Genre.Name
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	return &movieGenres, nil

}

func (m *DBModel) GetMovie(id int) (*Movie, error) {
	movie := &Movie{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT id, title, description, year, release_date, runtime, rating,mpaa_rating,created_at, updated_at FROM movies WHERE id = $1" //? doesn't work, use $1
	//QueryRowContext is a slow operation, so we use context to cancel it
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Year, &movie.ReleaseDate, &movie.Runtime, &movie.Rating, &movie.MPAARating, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	genres, err := getGenres(m, id, ctx)
	if err != nil {
		return nil, err
	}
	movie.MovieGenres = *genres

	return movie, nil

}
func (m *DBModel) GetAll() ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT id, title, description, year, release_date, runtime, rating,mpaa_rating,created_at, updated_at FROM movies ORDER BY title"
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)
	movies := make([]*Movie, 0)
	for rows.Next() {
		movie := &Movie{}

		err = rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Year, &movie.ReleaseDate, &movie.Runtime, &movie.Rating, &movie.MPAARating, &movie.CreatedAt, &movie.UpdatedAt)
		if err != nil {
			return nil, err
		}

		genres, err := getGenres(m, int(movie.ID), ctx)
		if err != nil {
			return nil, err
		}
		movie.MovieGenres = *genres

		movies = append(movies, movie)
	}
	return movies, nil
}
func (m *DBModel) GetGenres() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT id, genre_name FROM genres"
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)
	genres := make([]*Genre, 0)
	for rows.Next() {
		genre := &Genre{}

		err = rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}
	return genres, nil
}
