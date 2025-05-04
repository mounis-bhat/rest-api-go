package store

import (
	"database/sql"
	"time"

	"github.com/mounis-bhat/rest-api-go/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func (t *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.Insert(token)

	return token, err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	_, err := t.db.Exec("INSERT INTO tokens (user_id, hash, created_at, expiry, scope) VALUES ($1, $2, $3, $4, $5)",
		token.UserID,
		token.Hash,
		token.CreatedAt,
		token.Expiry,
		token.Scope,
	)

	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	_, err := t.db.Exec("DELETE FROM tokens WHERE user_id = $1 AND scope = $2", userID, scope)
	return err
}
