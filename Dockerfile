FROM golang:1.20-alpine AS builder
WORKDIR /app


COPY go.mod ./

RUN go mod download && go mod tidy


COPY . .

RUN ls -la /app


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o main ./main.go > /tmp/build.log 2>&1 || (cat /tmp/build.log && false)

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations/ ./migrations/

RUN apk update && apk add --no-cache ca-certificates

CMD ["./main"]