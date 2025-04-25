# Dockerfile

# Stage 1: Build
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o publisher .

# Stage 2: Run
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/publisher .

EXPOSE 8080
ENTRYPOINT ["./publisher"]
