package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Создаем новый фильм!")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	// в момент запроса данные будут в контексте роутера и их можно вытащить ParamsFromContext
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Отображение детальной информации о фильме с ид %d\n", id)

}
