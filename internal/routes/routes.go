package routes

import (
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

	return r
}
