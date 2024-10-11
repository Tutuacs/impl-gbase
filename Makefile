build:
	@go build -o bin/project-teste3 cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/project-teste3

migration:
	@go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down
