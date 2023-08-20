run:
	APP_ENV=local go run main.go

build:
	go build

unit-test:
	go test ./... -short -timeout 10s

mockgen:
	~/go/bin/mockgen -destination=internal/todo/mocks/mock_todo_repository.go -package mocks go-to-do-app/internal/todo Repository

db-test:
	go clean -testcache && go test ./internal/couchbase -v
