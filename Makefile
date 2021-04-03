.SILENT:

swagger-gen:
	swag init -g ./internal/app/app.go -o ./docs/swagger

vendor:
	go mod vendor

app-build: vendor swagger-gen
	go build -o ./bin -i ./cmd/app

app-dev: vendor
	go run ./cmd/app

mock-gen:
	go generate ./...

tests: vendor mock-gen
	go test ./... -v

dstop:
	docker-compose -f "docker-compose.yml" down

dstart: swagger-gen dstop
	docker-compose -f "docker-compose.yml" up -d --build
