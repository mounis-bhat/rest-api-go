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

	return r
}
