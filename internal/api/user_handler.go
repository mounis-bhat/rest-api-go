package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/mounis-bhat/rest-api-go/internal/store"
	"github.com/mounis-bhat/rest-api-go/internal/utils"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(store store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{userStore: store, logger: logger}
}

func (h *UserHandler) validateRegisterRequest(reg *registerUserRequest) error {
	if reg.Username == "" {
		return errors.New("username is required")
	}
	if len(reg.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(reg.Username) > 20 {
		return errors.New("username must be at most 20 characters long")
	}

	if reg.Email == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(reg.Email) {
		return errors.New("invalid email format")
	}
	if reg.Password == "" {
		return errors.New("password is required")
	}

	if len(reg.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(reg.Password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(reg.Password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(reg.Password)

	if !hasLower || !hasUpper || !hasDigit {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}
	if len(reg.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(reg.Password) > 20 {
		return errors.New("password must be at most 20 characters long")
	}

	return nil
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var reg registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	err = h.validateRegisterRequest(&reg)
	if err != nil {
		h.logger.Printf("Validation error: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: reg.Username,
		Email:    reg.Email,
	}

	err = user.PasswordHash.Set(reg.Password)
	if err != nil {
		h.logger.Printf("Error setting password hash: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to set password"})
		return
	}
	user, err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("Error creating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create user"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"user": user})

}

func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadIdParam(r)
	if err != nil {
		h.logger.Printf("Error reading user ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}

	var reg registerUserRequest
	err = json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		h.logger.Printf("Error decoding request body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}
	err = h.validateRegisterRequest(&reg)
	if err != nil {
		h.logger.Printf("Validation error: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user := &store.User{
		ID:       userId,
		Username: reg.Username,
		Email:    reg.Email,
	}

	err = user.PasswordHash.Set(reg.Password)
	if err != nil {
		h.logger.Printf("Error setting password hash: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to set password"})
		return
	}
	err = h.userStore.UpdateUser(user)
	if err != nil {
		h.logger.Printf("Error updating user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to update user"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}

func (h *UserHandler) HandleGetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		h.logger.Printf("Error: username query parameter is required")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "username query parameter is required"})
		return
	}
	user, err := h.userStore.GetUserByUsername(username)
	if err != nil {
		h.logger.Printf("Error retrieving user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve user"})
		return
	}
	if user == nil {
		h.logger.Printf("User not found")
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "User not found"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"user": user})
}

func (h *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadIdParam(r)
	if err != nil {
		h.logger.Printf("Error reading user ID: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid user ID"})
		return
	}

	err = h.userStore.DeleteUser(userId)
	if err != nil {
		h.logger.Printf("Error deleting user: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to delete user"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *UserHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userStore.GetAllUsers()
	if err != nil {
		h.logger.Printf("Error retrieving users: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to retrieve users"})
		return
	}

	if len(users) == 0 {
		h.logger.Printf("No users found")
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "No users found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"users": users})
}
