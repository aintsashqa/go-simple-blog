.SILENT:

swagger-gen:
	swag init -g ./internal/app/app.go -o ./api/swagger

depend:
	go mod tidy | go mod vendor

app-build: depend swagger-gen
	go build -o ./bin -i ./cmd/app

app-dev: depend swagger-gen
	go run ./cmd/app

mock-gen:
	go generate ./...

tests: depend mock-gen
	go test ./... -v

dstop:
	docker-compose -f "docker-compose.yml" down

dstart: swagger-gen dstop
	docker-compose -f "docker-compose.yml" up -d --build
