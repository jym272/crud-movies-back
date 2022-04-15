package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func getGenres(m *DBModel, movieId int64, ctx context.Context) (*map[int]string, error) {
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
func (m *DBModel) MovieExists(id int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	query := "SELECT id FROM movies WHERE id = $1"
	row := m.DB.QueryRowContext(ctx, query, id)
	var movieId int
	err := row.Scan(&movieId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *DBModel) GetMovie(id int64) (*Movie, error) {
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

		genres, err := getGenres(m, movie.ID, ctx)
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
	query := "SELECT id, genre_name FROM genres ORDER BY genre_name"
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

func (m *DBModel) GetMoviesByGenreWithID(genreID int64) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT movies.id, title, description, year, release_date, runtime, rating,mpaa_rating,movies.created_at, movies.updated_at FROM movies  INNER JOIN movies_genres ON movies.id = movies_genres.movie_id INNER JOIN genres ON movies_genres.genre_id = genres.id WHERE genre_id = $1 ORDER BY movies.title"
	rows, err := m.DB.QueryContext(ctx, query, genreID)
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

		genres, err := getGenres(m, movie.ID, ctx)
		if err != nil {
			return nil, err
		}
		movie.MovieGenres = *genres

		movies = append(movies, movie)
	}
	return movies, nil
}

func (m *DBModel) GetGenreNameByID(genreID int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT genre_name FROM genres WHERE id = $1"
	var genreName string
	err := m.DB.QueryRowContext(ctx, query, genreID).Scan(&genreName)
	if err != nil {
		return "", err
	}
	return genreName, nil
}

func (m *DBModel) UpdateMovie(id int64, movie *Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "UPDATE movies SET title = $1, description = $2, year = $3, release_date = $4, runtime = $5, rating = $6, mpaa_rating = $7, updated_at = $8 WHERE id = $9"
	_, err := m.DB.ExecContext(ctx, query, movie.Title, movie.Description, movie.Year, movie.ReleaseDate, movie.Runtime, movie.Rating, movie.MPAARating, time.Now(), id)
	if err != nil {
		return err
	}
	//delete all genres for this movie
	query = "DELETE FROM movies_genres WHERE movie_id = $1"
	_, err = m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	//add new genres
	for key := range movie.MovieGenres {
		query = "INSERT INTO movies_genres (movie_id, genre_id, created_at) VALUES ($1, $2, $3)"
		_, err = m.DB.ExecContext(ctx, query, id, key, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}
func (m *DBModel) InsertMovie(movie *Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "INSERT INTO movies (title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err := m.DB.ExecContext(ctx, query, movie.Title, movie.Description, movie.Year, movie.ReleaseDate, movie.Runtime, movie.Rating, movie.MPAARating, time.Now(), time.Now())
	if err != nil {
		return err
	}
	//get the new movie id
	query = "SELECT id FROM movies WHERE title = $1"
	var movieID int64
	err = m.DB.QueryRowContext(ctx, query, movie.Title).Scan(&movieID)
	if err != nil {
		return err
	}

	//add new genres
	for key := range movie.MovieGenres {
		query = "INSERT INTO movies_genres (movie_id, genre_id, created_at) VALUES ($1, $2, $3)"
		_, err = m.DB.ExecContext(ctx, query, movieID, key, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}
