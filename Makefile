# execution

.PHONY: build
build:
	go build -o ./bin/wallet ./cmd/wallet/main.go

start:
	./bin/wallet

run:
	go run ./cmd/wallet/main.go

# testing

test:
	go test ./... -cover -coverprofile=coverage.dev

.PHONY: coverage
coverage:
	go tool cover -html=coverage.dev
