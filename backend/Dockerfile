FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .

RUN swag init -g cmd/api/docs.go -o docs --parseDependency --parseInternal

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]
