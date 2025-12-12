build:
	@go build -o bin/crypto_exchange main.go


test:
	@go test -v ./...

run: build
	@./bin