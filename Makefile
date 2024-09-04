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

## vendor: очистка зависимостей и вендоринг зависимостей
.PHONY: vendor
vendor:
	@echo 'Проверка зависимостей и верификация модулей'
	go mod tidy
	go mod verify
	@echo 'Вендоринг зависимостей'
	go mod vendor


# ============================================================================= #
# BUILDING
# ============================================================================= #

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flag = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: сборка бинарника приложения
.PHONY: build/api
build/api:
	@echo 'Создание из cmd/api...'
	go build -ldflags=${linker_flag} -o=./bin/api ./cmd/api
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flag} -o=./bin/linux_amd64/api ./cmd/api
	

# ============================================================================= #
# PRODUCTION
# ============================================================================= #

production_host_ip = 'prod'

## production/connect: Подключиться к проду
.PHONY: production/connect
production/connect:
	ssh gl@${production_host_ip}


## production/deploy/api: Выгрузить на прод обновление
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api gl@${production_host_ip}:~
	rsync -rP --delete ./migrations gl@${production_host_ip}:~
	ssh -t gl@${production_host_ip} 'migrate -path ~/migrations -database $$GL_API_DSN up'