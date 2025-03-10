build:
	@go build -o bin/vinyl_store

run: build
	@./bin/vinyl_store

test:
	@go test -v ./...