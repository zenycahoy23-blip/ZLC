# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY schema.sql .
COPY init-db.sh .

RUN chmod +x init-db.sh

EXPOSE 8080

CMD ["./main"]
