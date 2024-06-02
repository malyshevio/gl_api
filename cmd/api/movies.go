package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gl_api.malyshev.io/internal/data"
	"gl_api.malyshev.io/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json."title"`
		Year    int32        `json."year"`
		Runtime data.Runtime `json."runtime"`
		Genres  []string     `json."genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validation section
	v := validator.New()

	v.Check(input.Title != "", "title", "Не может быть пустым")
	v.Check(len(input.Title) <= 500, "title", "Должно быть меньше 500 байт")

	v.Check(input.Year != 0, "year", "Не может быть пустым")
	v.Check(input.Year >= 1888, "year", "Год должен быть больше 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "Год не может быть из будующего")

	v.Check(input.Genres != nil, "genres", "Не может быть пустым")
	v.Check(len(input.Genres) >= 1, "genres", "должен быть заполнен хотябы 1 жанр")
	v.Check(len(input.Genres) <= 5, "genres", "Не более 5 жанров")

	v.Check(validator.Unique(input.Genres), "genres", "Жанры должны быть уникальными")

	v.Check(input.Runtime != 0, "runtime", "Не может быть пустым")
	v.Check(input.Runtime > 0, "runtime", "Должно быть положительным числом")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) unmarshalHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Foo string `json:"foo"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
