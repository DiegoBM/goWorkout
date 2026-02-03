package store

import (
	"database/sql"
	"time"

	"github.com/DiegoBM/goWorkout/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func (s *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = s.Insert(token)
	return token, err
}

func (s *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, scope, expiry)
	VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(query, token.Hash, token.UserID, token.Scope, token.Expiry)
	return err
}

func (s *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `DELETE FROM tokens WHERE scope = $1 AND user_id = $2`

	res, err := s.db.Exec(query, scope, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
