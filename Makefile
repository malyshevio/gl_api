## help: выведет помощь
help:
	@echo 'Использование:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	
## run/api: запуск API cmd/api
run/api:
	@go run ./cmd/api

## db/psql: подключение к базе консольным клиентом
db/psql:
	psql ${GL_API_DSN}

## db/migrations/up: запуск миграций на установку 
db/migrations/up: confirm
	@echo 'Запуск миграций...'
	migrate -path ./migrations -database ${GL_API_DSN} up

## db/migrations/new: создание новой пары миграции up/down 
db/migrations/new:
	@echo 'Создание пары миграций up/down ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## test/echo: как работает @
test/echo:
	echo 'ECHO'
	@echo 'muted ECHO'

confirm:
	@echo -n 'Вы уверены [y/N]' && read ans && [ $${ans:-N} = y ]
