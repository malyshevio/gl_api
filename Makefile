run:
	@go run ./cmd/api

psql:
	psql ${GL_API_DSN}

up: 
	@echo 'Запуск миграций...'
	migrate -path ./migrations -database ${GL_API_DSN} up

migration:
	@echo 'Создание пары миграций up/down ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

echo:
	echo 'ECHO'
	@echo 'muted ECHO'