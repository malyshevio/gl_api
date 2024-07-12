package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"gl_api.malyshev.io/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // timestamp when added to DB
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty,string"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "Не может быть пустым")
	v.Check(len(movie.Title) <= 500, "title", "Должно быть меньше 500 байт")

	v.Check(movie.Year != 0, "year", "Не может быть пустым")
	v.Check(movie.Year >= 1888, "year", "Год должен быть больше 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "Год не может быть из будующего")

	v.Check(movie.Genres != nil, "genres", "Не может быть пустым")
	v.Check(len(movie.Genres) >= 1, "genres", "должен быть заполнен хотябы 1 жанр")
	v.Check(len(movie.Genres) <= 5, "genres", "Не более 5 жанров")

	v.Check(validator.Unique(movie.Genres), "genres", "Жанры должны быть уникальными")

	v.Check(movie.Runtime != 0, "runtime", "Не может быть пустым")
	v.Check(movie.Runtime > 0, "runtime", "Должно быть положительным числом")
}

type MovieModel struct {
	DB *sql.DB
}

// Insert method to movie DB
func (m MovieModel) Insert(movie *Movie) error {
	query := `
		INSERT INTO movies (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
	`

	// TODO pattern to snippet storage
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Get method from movie DB
func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// TODO remove test pg_sleep(10)
	query := `
		SELECT pg_sleep(10), id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = $1
	`
	var movie Movie

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&[]byte{}, // pg_sleep()
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

// Update method to movie DB
func (m MovieModel) Update(movie *Movie) error {
	query := `
		UPDATE movies
		SET title=$1, year=$2, runtime=$3, genres=$4, version=version+1
		WHERE id = $5 AND version = $6
		RETURNING version
	`
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete from movie DB method
func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM movies
	WHERE id = $1
	`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// type MockMovieModel struct{}

// func (m MockMovieModel) Insert(movie *Movie) error {
// 	// mock
// }

// func (m MockMovieModel) Get(id int64) (*Movie, error) {
// 	// mock 2
// }

// func (m MockMovieModel) Update(movie *Movie) error {
// 	// mock 3
// }

// func (m MockMovieModel) Delete(id int64) error {
// 	// mock 4
// }
