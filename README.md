### Structure
```
├─ /bin
├─ /cmd
│　└─ /api
│　　　├─ healthcheck.go
│　　　├─ helpers.go
│　　　├─ main.go
│　　　├─ movies
│　　　└─ routers.go
├─ /internal
│　└─ /data
│　　　└─ movies.go
├─ /migration
├─ /remote
├─ go.mod
└─ Makefile
```

- `bin` собственно скомпилированный проект, подготовленный для деплоя например на прод
- `cmd/api` код самого приложения. Код для запуска сервера, чтения http запросов и аутентификации
- `internal` связанные вспомогательные пакеты для АПИ. Код взаимодействия с базой, валидация, отправка имейлов... Короче все что не специфично именно для *этого* приложения, всякое что можно переиспользовать где угодно еще. Отсюда будет *импортироваться* код в `cmd/api`, **НО** ни кода не наоборот!
- `migration` для SQL миграций тут ничего хитрого
-  `remote` конфиги и скрипты для продакшена
- `go.mod` собственно файл модуля
- `Makefile` билдер сценарий

### Endpoints

| Method | URL Pattern     | Handler            | Action                               |
| ------ | --------------- | ------------------ | ------------------------------------ |
| GET    | /v1/healthcheck | healthcheckHandler | Выведем немного информации о проекте |
| GET    | /v1/movies/:id  | showMovieHandler   | Показать детали конкретного фильма   |
| GET    | /v1/movies      | listMoviesHandler  | Отобразить все фильмы с фильтрами    |
| POST   | /v1/movies      | createMovieHandler | Создать новый фильм                  |
| PATCH  | /v1/movies/:id  | editMovieHandler   | Обновить информацию о фильме         |
| DELETE | /v1/movies/:id  | deleteMovieHandler | Удалить фильм из базы                |


## Migrations

тулинг для миграции
https://github.com/golang-migrate/migrate

```shell
$ migrate -path=./migrations -database=$GL_API_DSN up
```


## Фильтры
пример 1:

`/v1/movies?title=godzilla&genres=scifi,drama&page=1&page_size=5&sort=-year`


## Логи

| Ключ       | Описание                                        |
| ---------- | ----------------------------------------------- |
| level      | Уровень лога (INFO, ERROR, FATAL)               |
| time       | UTC время                                       |
| message    | сообщение в свободной форме                     |
| properties | дополнительные параметры например ключ\значения |
| trace      | стек вызова для отладки                         |

