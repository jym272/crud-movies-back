package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func (m *DBModel) GetMovie(id int) (*Movie, error) {
	movie := &Movie{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT id, title, description, year, release_date, runtime, rating,mpaa_rating,created_at, updated_at FROM movies WHERE id = $1" //? doesn't work, use $1
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Year, &movie.ReleaseDate, &movie.Runtime, &movie.Rating, &movie.MPAARating, &movie.CreatedAt, &movie.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return movie, nil
}
func (m *DBModel) GetAll() ([]*Movie, error) {
	//rows, err := m.DB.Query("SELECT id, title, year, director FROM movies")
	//if err != nil {
	//	return nil, err
	//}
	//defer rows.Close()
	//
	//movies := []*Movie{}
	//for rows.Next() {
	//	movie := &Movie{}
	//	err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Director)
	//	if err != nil {
	//		return nil, err
	//	}
	//	movies = append(movies, movie)
	//}
	//if err = rows.Err(); err != nil {
	//	return nil, err
	//}
	//return movies, nil
	return nil, nil
}
