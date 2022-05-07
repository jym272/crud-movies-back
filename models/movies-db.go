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

func (m *DBModel) MovieTitleExist(title string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	query := "SELECT id FROM movies WHERE title = $1"

	rows := m.DB.QueryRowContext(ctx, query, title)
	var id int64
	err := rows.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
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
	query := "SELECT id, title, description, year, release_date, runtime, rating,mpaa_rating,created_at, updated_at,coalesce(poster,'') FROM movies WHERE id = $1" //? doesn't work, use $1
	//QueryRowContext is a slow operation, so we use context to cancel it
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Year, &movie.ReleaseDate, &movie.Runtime, &movie.Rating, &movie.MPAARating, &movie.CreatedAt, &movie.UpdatedAt, &movie.Poster)
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
func (m *DBModel) GetAll(userID ...int64) ([]*Movie, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rows *sql.Rows
	var err error

	if len(userID) > 0 {
		userID_ := userID[0]
		query := "SELECT id, title, description, year, release_date, runtime, rating,mpaa_rating,created_at,  updated_at, coalesce(poster,'') FROM movies WHERE user_id = $1 ORDER BY title"
		rows, err = m.DB.QueryContext(ctx, query, userID_)

	} else {
		query := "SELECT id, title, description, year, release_date, runtime, rating,mpaa_rating,created_at,  updated_at, coalesce(poster,'') FROM movies ORDER BY title"
		rows, err = m.DB.QueryContext(ctx, query)
	}
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

		err = rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Year, &movie.ReleaseDate, &movie.Runtime, &movie.Rating, &movie.MPAARating, &movie.CreatedAt, &movie.UpdatedAt, &movie.Poster)
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
	query := "SELECT movies.id, title, description, year, release_date, runtime, rating,mpaa_rating,movies.created_at, movies.updated_at, coalesce(poster,'') FROM movies  INNER JOIN movies_genres ON movies.id = movies_genres.movie_id INNER JOIN genres ON movies_genres.genre_id = genres.id WHERE genre_id = $1 ORDER BY movies.title"
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

		err = rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Year, &movie.ReleaseDate, &movie.Runtime, &movie.Rating, &movie.MPAARating, &movie.CreatedAt, &movie.UpdatedAt, &movie.Poster)
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
	query := "UPDATE movies SET title = $1, description = $2, year = $3, release_date = $4, runtime = $5, rating = $6, mpaa_rating = $7, updated_at = $8, poster= $9 WHERE id = $10"
	_, err := m.DB.ExecContext(ctx, query, movie.Title, movie.Description, movie.Year, movie.ReleaseDate, movie.Runtime, movie.Rating, movie.MPAARating, time.Now(), movie.Poster, id)
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

func (m *DBModel) IsOwner(movieID int64, userID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT id FROM movies WHERE id = $1 AND user_id = $2"
	var id int64
	err := m.DB.QueryRowContext(ctx, query, movieID, userID).Scan(&id)
	if err != nil {
		return false, err
	}
	if id > 0 {
		return true, nil
	}
	return false, nil
}

func (m *DBModel) InsertMovie(movie *Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "INSERT INTO movies (title, description, year, release_date, runtime, rating, mpaa_rating, created_at, updated_at, poster, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	_, err := m.DB.ExecContext(ctx, query, movie.Title, movie.Description, movie.Year, movie.ReleaseDate, movie.Runtime, movie.Rating, movie.MPAARating, time.Now(), time.Now(), movie.Poster, movie.UserID)
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

func (m *DBModel) DeleteMovie(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//first delete from movies_genres
	query := "DELETE FROM movies_genres WHERE movie_id = $1"
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	//then delete from movies
	query = "DELETE FROM movies WHERE id = $1"
	_, err = m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) GetUser(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT id, username, password, created_at, updated_at FROM users WHERE username = $1"
	var user User
	err := m.DB.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (m *DBModel) CreateUser(user_ *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "INSERT INTO users (username, password, created_at, updated_at) VALUES ($1, $2, $3, $4)"
	_, err := m.DB.ExecContext(ctx, query, user_.Username, user_.Password, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}
func (m *DBModel) UpdateUserPasswordByUsername(username string, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "UPDATE users SET password = $1, updated_at = $2 WHERE username = $3"
	_, err := m.DB.ExecContext(ctx, query, password, time.Now(), username)
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) ValidateUser(userId int64, username string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT id, username, password, created_at, updated_at FROM users WHERE id = $1 AND username = $2"
	var user User
	err := m.DB.QueryRowContext(ctx, query, userId, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return false
	}
	return true
}

func (m *DBModel) IsFav(movieId int64, userId int64) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT id FROM favorite_movies WHERE movie_id = $1 AND user_id = $2"
	var id int64
	err := m.DB.QueryRowContext(ctx, query, movieId, userId).Scan(&id)
	if err != nil {
		return false
	}
	return true
}
func (m *DBModel) AddToFav(movieId int64, userId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "INSERT INTO favorite_movies (movie_id, user_id, created_at) VALUES ($1, $2, $3)"
	_, err := m.DB.ExecContext(ctx, query, movieId, userId, time.Now())
	if err != nil {
		return err
	}
	return nil
}
func (m *DBModel) RemoveFromFav(movieId int64, userId int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "DELETE FROM favorite_movies WHERE movie_id = $1 AND user_id = $2"
	_, err := m.DB.ExecContext(ctx, query, movieId, userId)
	if err != nil {
		return err
	}
	return nil
}
func (m *DBModel) GetFavorites(userId int64) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id, title, description, year, release_date, runtime, rating,mpaa_rating,created_at,  updated_at, coalesce(poster,'') FROM movies WHERE id IN (SELECT movie_id FROM favorite_movies WHERE user_id = $1) ORDER BY title"
	rows, err := m.DB.QueryContext(ctx, query, userId)

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

		err = rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Year, &movie.ReleaseDate, &movie.Runtime, &movie.Rating, &movie.MPAARating, &movie.CreatedAt, &movie.UpdatedAt, &movie.Poster)
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

func (m *DBModel) GetMoviesIds(genreID int64, search string, userID int64) ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var rows *sql.Rows
	var err error
	var query string
	switch {
	case search != "":
		query = "SELECT id FROM movies WHERE title ILIKE $1 ORDER BY title"
		rows, err = m.DB.QueryContext(ctx, query, "%"+search+"%")
	case genreID != 0:
		query = "SELECT id FROM movies WHERE id IN (SELECT movie_id FROM movies_genres WHERE genre_id = $1) ORDER BY title"
		rows, err = m.DB.QueryContext(ctx, query, genreID)
	case userID != 0:
		query = "SELECT id FROM movies WHERE id IN (SELECT movie_id FROM favorite_movies WHERE user_id = $1) ORDER BY title"
		rows, err = m.DB.QueryContext(ctx, query, userID) //m.User.ID
	default:
		query = "SELECT id FROM movies ORDER BY title"
		rows, err = m.DB.QueryContext(ctx, query)

	}

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)
	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
