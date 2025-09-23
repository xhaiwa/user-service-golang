# Étape de build
FROM golang:1.25.1-alpine AS builder
WORKDIR /app

# Copier les fichiers de modules Go
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste du code
COPY . .

# Compiler le binaire
RUN go build -o user-service ./cmd/main.go

# Étape runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user-service .
EXPOSE 8080

CMD ["./user-service"]