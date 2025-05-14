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

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(store store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{workoutStore: store, logger: logger}
}

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

func (h *WorkoutHandler) HandleGetAllWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.workoutStore.GetAllWorkouts()
	if err != nil {
		h.logger.Printf("Error retrieving workouts: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve workouts"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workouts": workouts})
}
