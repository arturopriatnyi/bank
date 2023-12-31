# execution

.PHONY: build
build:
	go build -o ./bin/bank ./cmd/bank/main.go

start:
	./bin/bank

run:
	docker-compose up --build mongodb postgresql & go run cmd/bank/main.go

# testing

test:
	go test ./... -cover -coverprofile=coverage.dev

.PHONY: coverage
coverage:
	go tool cover -html=coverage.dev
