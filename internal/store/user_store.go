package store

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.plaintext = &plaintext
	p.hash = hash
	return nil
}

func (p *password) Check(plaintext string) bool {
	if err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext)); err != nil {
		return false
	}
	return true
}

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(user *User) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int64) error
	GetAllUsers() ([]*User, error)
	GetUserToken(scope, tokenPlaintext string) (*User, error)
}

func (s *PostgresUserStore) GetUserToken(scope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `SELECT u.id, u.username, u.email, u.password_hash, u.created_at, u.updated_at
		FROM users u
		INNER JOIN tokens t ON u.id = t.user_id
		WHERE t.scope = $1 AND t.token_hash = $2 AND t.expires_at > $3`

	user := &User{
		PasswordHash: password{},
	}

	err := s.db.QueryRow(query, scope, tokenHash[:], time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if !user.PasswordHash.Check(tokenPlaintext) {
		return nil, fmt.Errorf("invalid token")
	}
	return user, nil
}

func (s *PostgresUserStore) CreateUser(user *User) (*User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err = tx.QueryRow(query, user.Username, user.Email, user.PasswordHash.hash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, tx.Commit()
}

func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `
  		SELECT id, username, email, password_hash, created_at, updated_at
  		FROM users
  		WHERE username = $1
  	`

	err := s.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	if user.ID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE users SET username = $1, email = $2, password_hash = $3,
		updated_at = NOW() WHERE id = $4`
	_, err = tx.Exec(query, user.Username, user.Email, user.PasswordHash.hash, user.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostgresUserStore) DeleteUser(id int64) error {
	if id == 0 {
		return fmt.Errorf("user ID is required")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM users WHERE id = $1`
	_, err = tx.Exec(query, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
func (s *PostgresUserStore) GetAllUsers() ([]*User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at
		FROM users`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		user := &User{
			PasswordHash: password{},
		}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
