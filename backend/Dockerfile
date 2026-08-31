FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /medical-booking .

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /medical-booking /app/medical-booking
COPY --from=builder /app/.env.docker /app/.env

EXPOSE 8080

CMD ["./medical-booking"]
