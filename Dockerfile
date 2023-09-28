FROM golang:1.19-alpine

WORKDIR ./bank
COPY . .

RUN go build -o ./bin/bank ./cmd/bank/main.go
CMD ["./bin/bank"]
