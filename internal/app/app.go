package app

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/mounis-bhat/rest-api-go/internal/api"
	"github.com/mounis-bhat/rest-api-go/internal/middleware"
	"github.com/mounis-bhat/rest-api-go/internal/store"
	"github.com/mounis-bhat/rest-api-go/internal/utils"
	"github.com/mounis-bhat/rest-api-go/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	Middleware     middleware.UserMiddleware
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	db, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		db.Close()
		return nil, err
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	workoutStore := store.NewPostgresWorkoutStore(db)
	userStore := store.NewPostgresUserStore(db)
	tokenStore := store.NewPostgresTokenStore(db)

	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(userStore, tokenStore, logger)
	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		Middleware:     middlewareHandler,
		DB:             db,
	}
	return app, nil
}

// HealthCheckHandler provides a health check endpoint
//
//	@Summary		Health check
//	@Description	Check if the API service is running and healthy
//	@Tags			Health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Service is healthy"
//	@Router			/health [get]
func (a *Application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"status": "OK"})
}
