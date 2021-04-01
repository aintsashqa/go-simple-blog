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

dstop:
	docker-compose -f "docker-compose.yml" down

dstart: dstop
	docker-compose -f "docker-compose.yml" up -d --build
