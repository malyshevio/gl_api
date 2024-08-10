package main

import (
	"context"
	"net/http"

	"gl_api.malyshev.io/internal/data"
)

// собственный тип на основе строки
type contextKey string

// конвертнем строку в тип contextKey и присвоит это константе - будем использовать эту константу как ключ контекста
const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("отсутствует значение для пользователя в контексте запроса")
	}

	return user
}
