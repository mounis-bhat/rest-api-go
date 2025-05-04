package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mounis-bhat/rest-api-go/internal/app"
)

func InitializeRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheckHandler)

	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkout)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkout)
	r.Get("/workouts", app.WorkoutHandler.HandleGetAllWorkouts)

	r.Post("/register", app.UserHandler.HandleCreateUser)
	r.Get("/user", app.UserHandler.HandleGetUserByUsername)
	r.Put("/users/{id}", app.UserHandler.HandleUpdateUser)
	r.Delete("/users/{id}", app.UserHandler.HandleDeleteUser)
	r.Get("/users", app.UserHandler.HandleGetAllUsers)

	r.Post("/tokens/auth", app.TokenHandler.HandleCreateToken)

	return r
}
