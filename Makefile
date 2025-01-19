build:
	@go build -o bin/maSrifiyy

run: build
	@./bin/maSrifiyy

test:
	@go test -v ./...
