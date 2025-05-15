# REST API GO

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/yourusername/rest-api-go/actions)
[![Go Version](https://img.shields.io/badge/go-1.24.2+-blue)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/license-MIT-yellow.svg)](LICENSE)

A robust REST API for managing workout routines built with Go, Chi router, and PostgreSQL.

*Last updated: May 14, 2025*

## Overview

This project provides a RESTful API for creating and retrieving workout routines, including detailed exercise entries. It's built using Go with the Chi router for HTTP routing and PostgreSQL for data persistence.

## Features

- Create workout routines with multiple exercise entries
- Retrieve workout details by ID
- Update existing workouts
- Delete workouts
- List all workouts
- Database migrations using Goose
- Modular architecture following clean code principles
- PostgreSQL database with transaction support

## Project Structure

```
rest-api-go/
├── docker-compose.yml    # Docker composition file
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── main.go               # Application entry point
├── internal/             # Internal application code
│   ├── api/              # API handlers
│   │   └── handler.go
│   ├── app/              # Application setup
│   │   └── app.go
│   ├── doc/              # API documentation
│   │   └── openapi.json
│   ├── middleware/       # HTTP middleware
│   │   └── middleware.go
│   ├── routes/           # HTTP routes
│   │   └── routes.go
│   ├── store/            # Database access
│   │   ├── database.go
│   │   ├── tokens.go
│   │   ├── feature_store.go
│   │   └── feature_store_test.go
│   ├── tokens/           # Token utilities
│   │   └── tokens.go
│   └── utils/            # Utility functions
│       └── utils.go
└── migrations/           # Database migrations
    ├── fs.go             # Embedded migrations
    └── migration_file.sql
```

## Prerequisites

- Go 1.24.2+
- PostgreSQL 17+ (containerized or standalone)

## Getting Started

### Setting up the database

Ensure you have a PostgreSQL instance running. You can use any PostgreSQL setup, whether containerized or standalone.

### Running the application

```bash
go run main.go
```

By default, the server runs on port 8080. You can specify a different port using the `-port` flag:

```bash
go run main.go -port 3000
```

## Documentation

API documentation is available in the OpenAPI format. You can find the specification in the `internal/doc/openapi.json` file.

## Development

### Database Migrations

This project uses [Goose v3](https://github.com/pressly/goose/v3) for database migrations. Migrations are automatically applied when the application starts using embedded SQL files.

### Adding New Endpoints

1. Create a new handler in the `internal/api` directory
2. Register the handler in the `internal/app/app.go` file
3. Add routes in the `internal/routes/routes.go` file
4. Update the OpenAPI documentation in `internal/doc/openapi.json`

## License

[MIT License](LICENSE)