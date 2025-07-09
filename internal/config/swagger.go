package config

import (
	"os"
)

type Config struct {
	Host     string
	Port     string
	BasePath string
}

func GetSwaggerConfig() Config {
	host := os.Getenv("SWAGGER_HOST")
	if host == "" {
		// Default to production host, can be overridden
		host = "workouts.mounis.net"
	}

	port := os.Getenv("SWAGGER_PORT")
	if port == "" {
		port = "443" // Default to HTTPS port for production
	}

	basePath := os.Getenv("SWAGGER_BASE_PATH")
	if basePath == "" {
		basePath = "/"
	}

	return Config{
		Host:     host,
		Port:     port,
		BasePath: basePath,
	}
}
