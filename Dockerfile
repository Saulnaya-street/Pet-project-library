FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod ./


RUN go mod download && go mod tidy

COPY . .

RUN ls -la /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o main ./db-service/cmd/main.go

FROM golang:1.22-alpine

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/db-service/migrations/ ./migrations/

RUN apk update && apk add --no-cache ca-certificates

CMD ["./main"]