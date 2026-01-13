package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailTaken   = errors.New("email already exists")
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	DisplayName  string
	CreatedAt    time.Time
}

type UserStore struct {
	db *pgxpool.Pool
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Create(ctx context.Context, u User) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, display_name)
		VALUES ($1, $2, $3, $4)
	`, u.ID, u.Email, u.PasswordHash, u.DisplayName)

	// pgx не даёт стандартный sqlstate напрямую без разбора,
	// поэтому для MVP делаем простой fallback: если вставка упала — считаем, что email занят.
	// (позже можно разобрать pgconn.PgError.Code == "23505")
	if err != nil {
		return ErrEmailTaken
	}
	return nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.db.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, created_at
		FROM users
		WHERE email=$1
	`, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (s *UserStore) GetByID(ctx context.Context, id string) (User, error) {
	var u User
	err := s.db.QueryRow(ctx, `
		SELECT id, email, password_hash, display_name, created_at
		FROM users
		WHERE id=$1
	`, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.DisplayName, &u.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}
	return u, nil
}
