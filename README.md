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
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── main.go               # Application entry point
├── internal/             # Internal application code
│   ├── api/              # API handlers
│   │   └── workout_handler.go
│   ├── app/              # Application setup
│   │   └── app.go
│   ├── routes/           # HTTP routes
│   │   └── routes.go
│   ├── store/            # Database access
│   │   ├── database.go
│   │   └── workout_store.go
│   └── utils/            # Utility functions
│       └── utils.go
└── migrations/           # Database migrations
    ├── fs.go             # Embedded migrations
    ├── 00001_users.sql
    ├── 00002_workouts.sql
    └── 00003_workout_entries.sql
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

## API Endpoints

### Health Check

```
GET /health
```

Example request:

```bash
curl http://localhost:8080/health
```

### User Management

#### Register a New User

```
POST /register
```

Example request:

```bash
curl -X POST \
  http://localhost:8080/register \
  -H 'Content-Type: application/json' \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "Password123"
  }'
```

#### Get User by Username

```
GET /user
```

Example request:

```bash
curl http://localhost:8080/user
```

#### Update a User

```
PUT /users/{id}
```

Example request:

```bash
curl -X PUT \
  http://localhost:8080/users/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "updateduser@example.com",
    "password": "NewPassword123"
  }'
```

#### Delete a User

```
DELETE /users/{id}
```

Example request:

```bash
curl -X DELETE http://localhost:8080/users/1
```

#### List All Users

```
GET /users
```

Example request:

```bash
curl http://localhost:8080/users
```

### Token Management

#### Create an Authentication Token

```
POST /tokens/auth
```

Example request:

```bash
curl -X POST \
  http://localhost:8080/tokens/auth \
  -H 'Content-Type: application/json' \
  -d '{
    "username": "existinguser",
    "password": "Password123"
  }'
```

### Workout Management

#### Create a New Workout

```
POST /workouts
```

Example request:

```bash
curl -X POST \
  http://localhost:8080/workouts \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Morning Cardio",
    "description": "Quick morning cardio routine",
    "duration_minutes": 30,
    "calories_burned": 250,
    "entries": [
      {
        "exercise_name": "Running",
        "sets": 1,
        "duration_seconds": 1200,
        "notes": "Moderate pace",
        "order_index": 0
      },
      {
        "exercise_name": "Jumping Jacks",
        "sets": 3,
        "reps": 20,
        "notes": "Full extension",
        "order_index": 1
      }
    ]
  }'
```

#### Get Workout by ID

```
GET /workouts/{id}
```

Example request:

```bash
curl http://localhost:8080/workouts/1
```

#### List All Workouts

```
GET /workouts
```

Example request:

```bash
curl http://localhost:8080/workouts
```

#### Update a Workout

```
PUT /workouts/{id}
```

Example request:

```bash
curl -X PUT \
  http://localhost:8080/workouts/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Updated Morning Cardio",
    "description": "Improved morning cardio routine",
    "duration_minutes": 45,
    "calories_burned": 300,
    "entries": [
      {
        "id": 1,
        "exercise_name": "Running",
        "sets": 1,
        "duration_seconds": 1500,
        "notes": "Faster pace",
        "order_index": 0
      },
      {
        "id": 2,
        "exercise_name": "Jumping Jacks",
        "sets": 4,
        "reps": 25,
        "notes": "Full extension",
        "order_index": 1
      }
    ]
  }'
```

#### Delete a Workout

```
DELETE /workouts/{id}
```

Example request:

```bash
curl -X DELETE http://localhost:8080/workouts/1
```

## Development

### Database Migrations

This project uses [Goose v3](https://github.com/pressly/goose/v3) for database migrations. Migrations are automatically applied when the application starts using embedded SQL files.

### Adding New Endpoints

1. Create a new handler in the `internal/api` directory
2. Register the handler in the `internal/app/app.go` file
3. Add routes in the `internal/routes/routes.go` file

## License

[MIT License](LICENSE)