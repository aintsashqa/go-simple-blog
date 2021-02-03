.SILENT:

swagger-gen:
	swag init -g ./internal/app/app.go -o ./docs/swagger

app-build: swagger-gen
	go build -o ./bin -i ./cmd/app

app-dev:
	go run ./cmd/app

mock-gen:
	go generate ./...

tests: mock-gen
	go test ./... -v
