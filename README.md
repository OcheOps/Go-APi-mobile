# Field Notes API

A production-style Go REST API for the "Field Notes" mobile assessment.

## Features
- **Language**: Go
- **Database**: SQLite (pure Go driver)
- **Router**: `go-chi`
- **Docs**: Swagger UI included
- **Deployment**: Dockerized

## Quick Start

### 1. Via Docker (Recommended)
```bash
# Build
docker build -t fieldnotes .

# Run (with data persistence)
docker run -d -p 8080:8080 -v $(pwd)/data:/app fieldnotes
```
*Note: The app writes `notes.db` to the working directory.*

### 2. Manual Run
```bash
go mod tidy
go run main.go db.go handlers.go models.go
# or
go build -o fieldnotes
./fieldnotes
```

## Documentation
Once running, access the Swagger UI at:
**[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

## Endpoints
- **GET** `/v1/notes`
- **POST** `/v1/notes`
- **PATCH** `/v1/notes/{id}`
- **DELETE** `/v1/notes/{id}`

See Swagger UI for details and test payloads.
