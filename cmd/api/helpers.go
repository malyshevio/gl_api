package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	// в момент запроса данные будут в контексте роутера и их можно вытащить ParamsFromContext
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("не правильный id")
	}

	return id, nil
}

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		// use errors.As for general type of errors
		case errors.As(err, &syntaxError):
			return fmt.Errorf("Некоректрый формат JSON %d", syntaxError.Offset)

			// use errors.Is for specific error
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("В теле запроса есть ошибки форматирования")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("тело запроса содержит некорректный JSON тип поля %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("тело запроса содержит некорректный JSON тип с символа %d", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("тело запроса не должно быть пустым")

			// TODO panicking vs erroring is discussable
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}

	}

	return nil

}
