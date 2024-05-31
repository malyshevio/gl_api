package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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
	// 1: ограничим размер тела запроса
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// 2: запрещаем лишние поля
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
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

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("тело запроса содержит неизвестное поле %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("тело запроса должно быть меньше %d байт", maxBytes)
			// panicking vs erroring is discussable

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}

	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("тело запроса должно содержать только 1 JSON")

	}

	return nil

}
