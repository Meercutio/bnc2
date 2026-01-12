package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID           string
	Email        string
	PasswordHash string
	DisplayName  string
}

type UserStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, u User) error {
	_, err := s.db.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, display_name)
		 VALUES ($1, $2, $3, $4)`,
		u.ID, u.Email, u.PasswordHash, u.DisplayName,
	)
	return err
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.db.QueryRow(ctx,
		`SELECT id, email, password_hash, display_name
		 FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName)

	if err != nil {
		return User{}, ErrUserNotFound
	}
	return u, nil
}

func (s *UserStore) GetByID(ctx context.Context, id string) (User, error) {
	var u User
	err := s.db.QueryRow(ctx,
		`SELECT id, email, password_hash, display_name
		 FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName)

	if err != nil {
		return User{}, ErrUserNotFound
	}
	return u, nil
}
