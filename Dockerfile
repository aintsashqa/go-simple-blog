FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o ./bin -i ./cmd/app

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist
RUN mkdir -p configs
RUN mkdir -p migrations

# Copy binary from build to main folder
RUN cp /build/bin/app .
RUN cp /build/configs/* ./configs
RUN cp /build/migrations/* ./migrations
RUN cp /build/.env .

# Command to run when starting the container
EXPOSE ${APP_PORT}
CMD ["/dist/app"]
