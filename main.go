// Package main REST API for workout tracking
//
//	@title			Workout Tracker API
//	@version		1.0
//	@description	A REST API for tracking workouts and managing users.
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host		localhost:8080
//	@BasePath	/
//	@schemes	http
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/mounis-bhat/rest-api-go/internal/app"
	"github.com/mounis-bhat/rest-api-go/internal/routes"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	var port int

	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	r := routes.InitializeRoutes(app)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://workouts.mounis.net"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Println("Starting application on port", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
