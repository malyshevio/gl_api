package data

import (
	"time"

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
