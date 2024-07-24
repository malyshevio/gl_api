package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoveryPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// дефер после паники
		defer func() {
			if err := recover(); err != nil {
				// в случае паники кинем хеадер "Connection: close" в ответ, это схлопнет  соединение
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
