# REST API GO

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/yourusername/rest-api-go/actions)
[![Go Version](https://img.shields.io/badge/go-1.24.2+-blue)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/license-MIT-yellow.svg)](LICENSE)

A robust REST API for managing workout routines built with Go, Chi router, and PostgreSQL.

_Last updated: May 14, 2025_

## Overview

This project provides a RESTful API for creating and retrieving workout routines, including detailed exercise entries. It's built using Go with the Chi router for HTTP routing and PostgreSQL for data persistence.

## Features

- Create workout routines with multiple exercise entries
- Retrieve workout details by ID
- Update existing workouts
- Delete workouts
- List all workouts
- User authentication and authorization with JWT tokens
- User registration and management
- **Modern API Documentation** - Clean, interactive OpenAPI docs powered by Scalar
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
├── .env.example          # Environment variables template
├── docs/                 # Auto-generated Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/             # Internal application code
│   ├── api/              # API handlers
│   │   ├── token_handler.go
│   │   ├── user_handler.go
│   │   └── workout_handler.go
│   ├── app/              # Application setup
│   │   └── app.go
│   ├── middleware/       # HTTP middleware
│   │   └── middleware.go
│   ├── routes/           # HTTP routes
│   │   └── routes.go
│   ├── store/            # Database access
│   │   ├── database.go
│   │   ├── tokens.go
│   │   ├── user_store.go
│   │   └── workout_store.go
│   ├── tokens/           # Token utilities
│   │   └── tokens.go
│   └── utils/            # Utility functions
│       └── utils.go
└── migrations/           # Database migrations
    ├── fs.go             # Embedded migrations
    └── *.sql             # Migration files
```

## Prerequisites

- Go 1.24.2+
- PostgreSQL 17+ (containerized or standalone)

## Getting Started

### Environment Setup

1. **Copy the environment template:**

   ```bash
   cp .env.example .env
   ```

2. **Configure your environment variables in `.env`:**

   ```bash
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=workout_tracker
   DB_USER=postgres
   DB_PASSWORD=your_password_here
   DB_SSL_MODE=disable

   # Application Configuration
   APP_PORT=8080
   JWT_SECRET=your_jwt_secret_here

   # Swagger Configuration (Optional - defaults to production values)
   SWAGGER_HOST=localhost:8080  # For local development
   # SWAGGER_HOST=workouts.mounis.net  # For production
   ```

### Setting up the database

Ensure you have a PostgreSQL instance running. You can use the provided Docker Compose file:

```bash
docker-compose up -d postgres
```

Or use any PostgreSQL setup, whether containerized or standalone.

### Running the application

```bash
go run main.go
```

By default, the server runs on port 8080. You can specify a different port using the `-port` flag:

```bash
go run main.go -port 3000
```

### Accessing the API Documentation

Once the server is running, you can access the interactive API documentation at:

```
http://localhost:8080/docs
```

## Documentation

### API Documentation

The API includes comprehensive OpenAPI documentation with a modern, clean interface powered by Scalar.

**Access the API Documentation:**

- Start the server: `go run main.go`
- Visit: `http://localhost:8080/docs`

**Available Documentation:**

- **Scalar UI**: Modern, clean documentation interface at `/docs`
- **JSON Spec**: Raw OpenAPI spec at `/swagger/doc.json`
- **YAML Spec**: Raw OpenAPI spec available in `docs/swagger.yaml`

**Generating Documentation:**
If you modify the API annotations, regenerate the documentation:

```bash
# Install swag CLI (if not already installed)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init
```

**Note on OpenAPI Version:**
Currently using Swagger 2.0 specification generated by the `swag` tool. While the specification format is 2.0, Scalar provides a much more modern and clean documentation interface compared to traditional Swagger UI. For future OpenAPI 3.0+ migration, consider using tools like `oapi-codegen` or manual specification creation.

### API Endpoints

#### Authentication

- `POST /register` - Register a new user
- `POST /tokens/auth` - Authenticate and get JWT token

#### Users (Protected)

- `GET /user` - Get user by username (query parameter)
- `GET /users` - Get all users
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

#### Workouts (Protected)

- `GET /workouts` - Get all workouts
- `GET /workouts/{id}` - Get workout by ID
- `POST /workouts` - Create new workout
- `PUT /workouts/{id}` - Update workout
- `DELETE /workouts/{id}` - Delete workout

#### Health

- `GET /health` - Health check endpoint

**Authentication:** Most endpoints require a Bearer token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Development

### Database Migrations

This project uses [Goose v3](https://github.com/pressly/goose/v3) for database migrations. Migrations are automatically applied when the application starts using embedded SQL files.

### Adding New Endpoints

1. Create a new handler in the `internal/api` directory
2. Add Swagger annotations to your handler functions using the following format:
   ```go
   // HandlerFunction does something
   //
   //	@Summary		Brief description
   //	@Description	Detailed description
   //	@Tags			TagName
   //	@Accept			json
   //	@Produce		json
   //	@Security		BearerAuth (if authentication required)
   //	@Param			paramName	path/query/body	type	required	"Description"
   //	@Success		200			{object}		ResponseType
   //	@Failure		400			{object}		ErrorResponse
   //	@Router			/endpoint [method]
   func (h *Handler) HandlerFunction(w http.ResponseWriter, r *http.Request) {
       // implementation
   }
   ```
3. Register the handler in the `internal/app/app.go` file
4. Add routes in the `internal/routes/routes.go` file
5. Regenerate documentation: `swag init`

### Code Style and Swagger Annotations

- Always include Swagger annotations for public API endpoints
- Use consistent response structures (`ErrorResponse`, custom response types)
- Include authentication requirements with `@Security BearerAuth`
- Provide examples in your struct tags:
  ```go
  type UserRequest struct {
      Username string `json:"username" example:"johndoe" validate:"required"`
      Email    string `json:"email" example:"user@example.com" validate:"required,email"`
  }
  ```

## Deployment

### Production Deployment

The project includes a GitHub Actions workflow (`.github/workflows/deploy.yml`) that automatically:

1. **Builds the application** with Go 1.24.0
2. **Generates Swagger documentation** for production host
3. **Deploys to VPS** via SCP
4. **Restarts the service** via SSH

**Important for Production:**

- The API documentation will be available at: `https://workouts-api.mounis.net/docs`
- The API host is configured for `workouts-api.mounis.net` in production
- CORS is configured for `https://workouts-api.mounis.net`

**Required GitHub Secrets:**

- `DATABASE_URL`: PostgreSQL connection string
- `HOST`: VPS IP address or hostname
- `USERNAME`: SSH username for VPS
- `PRIVATE_KEY`: SSH private key
- `PASSPHRASE`: SSH key passphrase (if applicable)

### Local Development vs Production

| Environment    | API Documentation              | CORS Origin                       | Port |
| -------------- | ------------------------------ | --------------------------------- | ---- |
| **Local**      | `localhost:8080/docs`          | Any                               | 8080 |
| **Production** | `workouts-api.mounis.net/docs` | `https://workouts-api.mounis.net` | 8080 |

## License

[MIT License](LICENSE)
