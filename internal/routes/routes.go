package routes

import (
	"encoding/json"
	"net/http"
	"os"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/go-chi/chi/v5"
	"github.com/mounis-bhat/rest-api-go/internal/app"
)

func InitializeRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		r.Get("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkoutByID))
		r.Post("/workouts", app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkout))
		r.Delete("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkout))
		r.Get("/workouts", app.Middleware.RequireUser(app.WorkoutHandler.HandleGetAllWorkouts))

		r.Get("/user", app.Middleware.RequireUser(app.UserHandler.HandleGetUserByUsername))
		r.Put("/users/{id}", app.Middleware.RequireUser(app.UserHandler.HandleUpdateUser))
		r.Delete("/users/{id}", app.Middleware.RequireUser(app.UserHandler.HandleDeleteUser))
		r.Get("/users", app.Middleware.RequireUser(app.UserHandler.HandleGetAllUsers))
	})

	r.Get("/health", app.HealthCheckHandler)
	r.Post("/register", app.UserHandler.HandleCreateUser)
	r.Post("/tokens/auth", app.TokenHandler.HandleCreateToken)

	// API Documentation with Scalar
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		// Read the swagger.json file content
		specContent, err := os.ReadFile("./docs/swagger.json")
		if err != nil {
			http.Error(w, "Failed to read API specification", http.StatusInternalServerError)
			return
		}

		// Parse JSON to ensure it's valid
		var spec map[string]interface{}
		if err := json.Unmarshal(specContent, &spec); err != nil {
			http.Error(w, "Invalid API specification", http.StatusInternalServerError)
			return
		}

		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecContent: spec,
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Workout Tracker API Documentation",
			},
			DarkMode: true,
		})
		if err != nil {
			http.Error(w, "Failed to generate documentation", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	})

	// Serve the JSON spec file
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	return r
}
