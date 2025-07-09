package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mounis-bhat/rest-api-go/internal/store"
	"github.com/mounis-bhat/rest-api-go/internal/tokens"
	"github.com/mounis-bhat/rest-api-go/internal/utils"
)

type TokenHandler struct {
	userStore  store.UserStore
	tokenStore store.TokenStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username" example:"johndoe" validate:"required"`       // Username for authentication
	Password string `json:"password" example:"SecurePass123" validate:"required"` // Password for authentication
}

type TokenResponse struct {
	AuthToken string `json:"auth_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT authentication token
}

func NewTokenHandler(userStore store.UserStore, tokenStore store.TokenStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		userStore:  userStore,
		tokenStore: tokenStore,
		logger:     logger,
	}
}

// HandleCreateToken authenticates a user and returns an auth token
//
//	@Summary		Authenticate user
//	@Description	Authenticate a user with username and password and return an auth token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		createTokenRequest	true	"User credentials"
//	@Success		200			{object}	TokenResponse		"Authentication successful"
//	@Failure		400			{object}	ErrorResponse		"Invalid request payload"
//	@Failure		401			{object}	ErrorResponse		"Invalid username or password"
//	@Failure		500			{object}	ErrorResponse		"Internal server error"
//	@Router			/tokens/auth [post]
func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Println("Error decoding request body:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload"})
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)

	if err != nil {
		h.logger.Println("Error fetching user:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	if user == nil {
		h.logger.Println("User not found")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid username or password"})
		return
	}

	if !user.PasswordHash.Check(req.Password) {
		h.logger.Println("Invalid password")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid username or password"})
		return
	}
	token, err := h.tokenStore.CreateNewToken(int(user.ID), 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Println("Error creating token:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"auth_token": token})
}
