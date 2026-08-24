# Boilerplate User CRUD Makefile

-include conf/config.env
export

run:
	@go run app/main.go

test:
	@go test -v ./...
