package main

import (
	"fmt"
	"net/http"
)

// lopError() общий метод хелпер для логирования сообщений, позже заменю на структурный логер
func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

// errorResponse() общий метод для отправки сообщений в формате JSON
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// serverErrorResponse() если в приложении все совсем плохо
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "на сервере проблемы и нет возможности обработать ваш запрос"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse() шлем 404
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "Запрошенный ресурс не найден"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse() шлем 405 когда нет слушателя на ресурсе
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("%s метод не поддерживается для этого ресурса", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "невозможно обновить запись, конфликт редактирования или запись удалена"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "Превышен лимит запросов"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "недействительные данные пользователя"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "отсутствует или неверный токен аутентификации"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "необходибо аутентификация для доступа к этому ресурсу"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "ваш аккаунт должен быть активирован для доступа к этому ресурсу"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "недостаточно привилегий для доступа"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
