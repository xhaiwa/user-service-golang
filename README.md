# User Service (Go + Gin)

## Description

Backend service to manage users and their preferences.

- Language: Go
- Framework: Gin
- Database: PostgreSQL
- Authentication: JWT / OAuth2 (planned)
- ORM: GORM (planned)
- Fully Dockerized for easy deployment and testing

## Prerequisites

- Docker >= 24.x
- Docker Compose >= 2.x

## 1. Build the Docker Image

```bash
# From the project root
docker build -t user-service:latest .
```

## 2. Run the Service with Docker

```bash
# Run the service alone
docker run -p 8080:8080 user-service:latest
```

- The service will be available at: <http://localhost:8080>

## 3. Run with Docker Compose (includes PostgreSQL)

```bash
docker-compose up -d
```

- User Service: <http://localhost:8080>

- PostgreSQL: port 5432 (configurable in docker-compose.yml)

## Health Check

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

## 4. Stop the Containers

```bash
docker-compose down
```

Notes

- Environment variables (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME) are defined in docker-compose.yml.

- The service is ready to be extended with JWT, OAuth2, and user CRUD endpoints.

