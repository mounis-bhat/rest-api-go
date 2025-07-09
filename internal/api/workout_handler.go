package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/mounis-bhat/rest-api-go/internal/middleware"
	"github.com/mounis-bhat/rest-api-go/internal/store"
	"github.com/mounis-bhat/rest-api-go/internal/utils"
)

type WorkoutEntryResponse struct {
	ID              int      `json:"id" example:"1"`                       // Entry ID
	ExerciseName    string   `json:"exercise_name" example:"Push ups"`     // Name of the exercise
	Sets            int      `json:"sets" example:"3"`                     // Number of sets
	Reps            *int     `json:"reps" example:"15"`                    // Number of repetitions
	DurationSeconds *int     `json:"duration_seconds" example:"60"`        // Duration in seconds
	Weight          *float64 `json:"weight" example:"75.5"`                // Weight in kg
	Notes           string   `json:"notes" example:"Good form maintained"` // Additional notes
	OrderIndex      int      `json:"order_index" example:"1"`              // Order of exercise in workout
}

type WorkoutResponse struct {
	ID              int                    `json:"id" example:"1"`                              // Workout ID
	UserID          int64                  `json:"user_id" example:"1"`                         // User ID who owns the workout
	Title           string                 `json:"title" example:"Morning Cardio"`              // Workout title
	Description     string                 `json:"description" example:"High intensity cardio"` // Workout description
	DurationMinutes int                    `json:"duration_minutes" example:"45"`               // Duration in minutes
	CaloriesBurned  int                    `json:"calories_burned" example:"350"`               // Calories burned
	CreatedAt       string                 `json:"created_at" example:"2024-01-01T12:00:00Z"`   // Creation timestamp
	UpdatedAt       string                 `json:"updated_at" example:"2024-01-01T12:00:00Z"`   // Last update timestamp
	Entries         []WorkoutEntryResponse `json:"entries"`                                     // List of workout exercises
}

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(store store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{workoutStore: store, logger: logger}
}

// HandleGetWorkoutByID retrieves a specific workout by ID
//
//	@Summary		Get workout by ID
//	@Description	Retrieve a specific workout and its exercises by workout ID
//	@Tags			Workouts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int				true	"Workout ID"
//	@Success		200	{object}	WorkoutResponse	"Workout details"
//	@Failure		400	{object}	ErrorResponse	"Invalid workout ID"
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		404	{object}	ErrorResponse	"Workout not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/workouts/{id} [get]
func (h *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		h.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout ID"})
		return
	}

	workout, err := h.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		h.logger.Printf("Error retrieving workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve workout"})
		return
	}

	if workout == nil {
		h.logger.Printf("Workout with ID %d not found", workoutId)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

// HandleCreateWorkout creates a new workout
//
//	@Summary		Create a new workout
//	@Description	Create a new workout with exercises for the authenticated user
//	@Tags			Workouts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			workout	body		store.Workout	true	"Workout data"
//	@Success		201		{object}	WorkoutResponse	"Workout created successfully"
//	@Failure		400		{object}	ErrorResponse	"Invalid request payload"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/workouts [post]
func (h *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		h.logger.Printf("Unauthorized user")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"})
		return
	}

	workout.UserID = currentUser.ID

	result, err := h.workoutStore.CreateWorkout(&workout)
	if err != nil {
		h.logger.Printf("Error creating workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create workout"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": result})
}

// HandleUpdateWorkout updates an existing workout
//
//	@Summary		Update workout
//	@Description	Update an existing workout and its exercises (only by the owner)
//	@Tags			Workouts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int				true	"Workout ID"
//	@Param			workout	body		store.Workout	true	"Updated workout data"
//	@Success		200		{object}	WorkoutResponse	"Workout updated successfully"
//	@Failure		400		{object}	ErrorResponse	"Invalid request data"
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		403		{object}	ErrorResponse	"Forbidden - not the owner"
//	@Failure		404		{object}	ErrorResponse	"Workout not found"
//	@Failure		500		{object}	ErrorResponse	"Internal server error"
//	@Router			/workouts/{id} [put]
func (h *WorkoutHandler) HandleUpdateWorkout(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		h.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout ID"})
		return
	}

	var workout store.Workout
	err = json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	workout.ID = int(workoutId)

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		h.logger.Printf("Unauthorized user")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"})
		return
	}

	workoutOwner, err := h.workoutStore.GetWorkoutOwner(workoutId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Printf("Workout with ID %d not found", workoutId)
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
			return
		}
		h.logger.Printf("Error retrieving workout owner: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve workout owner"})
		return
	}

	if int64(workoutOwner) != currentUser.ID {
		h.logger.Printf("User %d is not authorized to update workout %d", currentUser.ID, workoutId)
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "Forbidden"})
		return
	}

	err = h.workoutStore.UpdateWorkout(&workout)
	if err != nil {
		h.logger.Printf("Error updating workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update workout"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

// HandleDeleteWorkout deletes a workout
//
//	@Summary		Delete workout
//	@Description	Delete a workout by ID (only by the owner)
//	@Tags			Workouts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Workout ID"
//	@Success		204	"Workout deleted successfully"
//	@Failure		400	{object}	ErrorResponse	"Invalid workout ID"
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		403	{object}	ErrorResponse	"Forbidden - not the owner"
//	@Failure		404	{object}	ErrorResponse	"Workout not found"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/workouts/{id} [delete]
func (h *WorkoutHandler) HandleDeleteWorkout(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIdParam(r)
	if err != nil {
		h.logger.Printf("Error reading workout ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid workout ID"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		h.logger.Printf("Unauthorized user")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Unauthorized"})
		return
	}

	workoutOwner, err := h.workoutStore.GetWorkoutOwner(workoutId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Printf("Workout with ID %d not found", workoutId)
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
			return
		}
		h.logger.Printf("Error retrieving workout owner: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve workout owner"})
		return
	}

	if int64(workoutOwner) != currentUser.ID {
		h.logger.Printf("User %d is not authorized to delete workout %d", currentUser.ID, workoutId)
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "Forbidden"})
		return
	}

	err = h.workoutStore.DeleteWorkout(workoutId)
	if err != nil {
		if err == sql.ErrNoRows {
			h.logger.Printf("Workout with ID %d not found for deletion", workoutId)
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "Workout not found"})
			return
		}
		h.logger.Printf("Error deleting workout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to delete workout"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleGetAllWorkouts retrieves all workouts
//
//	@Summary		Get all workouts
//	@Description	Retrieve a list of all workouts in the system
//	@Tags			Workouts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		WorkoutResponse	"List of workouts"
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	ErrorResponse	"Internal server error"
//	@Router			/workouts [get]
func (h *WorkoutHandler) HandleGetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.workoutStore.GetAllWorkouts()
	if err != nil {
		h.logger.Printf("Error retrieving workouts: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve workouts"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workouts": workouts})
}
