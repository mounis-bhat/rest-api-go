package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mounis-bhat/rest-api-go/internal/store"
	"github.com/mounis-bhat/rest-api-go/internal/tokens"
	"github.com/mounis-bhat/rest-api-go/internal/utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}
func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("user not found in context")
	}
	return user
}

func (m *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(("Vary"), "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.SplitN(authHeader, " ", 2)
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
				"error": "Invalid authorization header",
			})
			return
		}

		token := headerParts[1]

		user, err := m.UserStore.GetUserToken(tokens.ScopeAuth, token)

		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{
				"error": "Internal server error",
			})
			return
		}

		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{
				"error": "Invalid or expired token",
			})
			return
		}

		SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
