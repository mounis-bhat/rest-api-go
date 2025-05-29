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
