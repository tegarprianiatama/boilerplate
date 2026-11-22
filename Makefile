# Boilerplate User CRUD Makefile

-include conf/config.env
export

run:
	@go run app/main.go

docker-up:
	@docker compose up -d

docker-down:
	@docker compose down

test:
	@go test -v ./...
