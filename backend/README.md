# TatApps Backend

Backend API untuk TatApps menggunakan Golang, Gin, dan PostgreSQL.

## Quick Start

```bash
# Install dependencies
go mod download

# Setup environment
cp .env.example .env
# Edit .env sesuai konfigurasi

# Run server
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

## API Documentation

Lihat [README.md](../README.md) untuk dokumentasi lengkap API endpoints.

## Development

```bash
# Run with auto-reload (install air first)
go install github.com/cosmtrek/air@latest
air
```

## Build for Production

```bash
go build -o tatapps-api cmd/api/main.go
./tatapps-api
```
