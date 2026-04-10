package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	return insertUser(ctx, s.db, u)
}

func (s *UserStore) CreateWithStats(ctx context.Context, u User, initialRating int) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := insertUser(ctx, tx, u); err != nil {
		return err
	}

	// Insert only the PK and rely on column defaults so registration works
	// against both pre-rating and current schemas.
	if _, err := tx.Exec(ctx, `
		INSERT INTO player_stats (user_id)
		VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
	`, u.ID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

type dbExec interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func insertUser(ctx context.Context, db dbExec, u User) error {
	_, err := db.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, display_name)
		VALUES ($1, $2, $3, $4)
	`, u.ID, u.Email, u.PasswordHash, u.DisplayName)

	if isUniqueViolation(err) {
		return ErrEmailTaken
	}
	return err
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
