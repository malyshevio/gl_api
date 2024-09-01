include .envrc

# ============================================================================= #
# HELPERS
# ============================================================================= #


## help: выведет помощь
.PHONY: help
help:
	@echo 'Использование:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Вы уверены [y/N]' && read ans && [ $${ans:-N} = y ]


# ============================================================================= #
# DEVELOPMENT
# ============================================================================= #

## run/api: запуск API cmd/api
.PHONY: run/api
run/api:
	@go run ./cmd/api -db-dsn=${GL_API_DSN}


## db/psql: подключение к базе консольным клиентом
.PHONY: db/psql
db/psql:
	psql ${GL_API_DSN}

## db/migrations/up: запуск миграций на установку 
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Запуск миграций...'
	migrate -path ./migrations -database ${GL_API_DSN} up

## db/migrations/new: создание новой пары миграции up/down 
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Создание пары миграций up/down ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}


# ============================================================================= #
# TESTS
# ============================================================================= #

## test/echo: как работает @
.PHONY: test/echo
test/echo:
	echo 'unmuted ECHO'
	@echo 'muted ECHO'

## test/echo: как работает @
.PHONY: test/env
test/env:
	@echo ${TEST_ECHO}


# ============================================================================= #
# QUALITY CONTROL
# ============================================================================= #

## audit: актуализация зависимостей, форматирование и тестирование кода. Запускать перед комитами в мастер
.PHONY: audit
audit:
	@echo 'Проверка зависимостей и верификация модулей'
	go mod tidy
	go mod verify
	@echo 'форматирование кода'
	go fmt ./...
	@echo 'Проверка кода vet'
	go vet ./...
	staticcheck ./...
	@echo 'Запуск тестов'
	go test -race -vet=off ./...