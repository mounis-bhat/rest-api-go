package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mounis-bhat/rest-api-go/internal/store"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
}

func NewWorkoutHandler(store store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{workoutStore: store}
}

func (h *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	workoutId, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	workout, err := h.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		http.Error(w, "Failed to retrieve workout", http.StatusInternalServerError)
		return
	}

	if workout == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
	w.WriteHeader(http.StatusOK)

}

func (h *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result, err := h.workoutStore.CreateWorkout(&workout)
	if err != nil {
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}
	workout = *result
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
	w.WriteHeader(http.StatusCreated)
}
