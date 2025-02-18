package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/codepnw/react_go_ecom/internal/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByID(ctx context.Context, id int) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	SaveRefreshToken(userID int, token string, expires time.Time) error
	ValidateRefreshToken(token string) (int, error)
	DeleteRefreshToken(token string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (email, password, role, enabled, created_at, updated_at, picture, address)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	var id int

	err := r.db.QueryRowContext(
		ctx,
		query,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Enabled,
		&user.CreatedAt,
		&user.UpdatedAt,
		"later",
		"later",
	).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	query := `
		SELECT id, email, password, role
		FROM users WHERE id = $1
	`
	var user entities.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, email, password, role
		FROM users WHERE email = $1
	`
	var user entities.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) SaveRefreshToken(userID int, token string, expires time.Time) error {
	query := `INSERT INTO refresh_token (user_id, token, expire_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, userID, token, expires)

	return err
}

func (r *userRepository) ValidateRefreshToken(token string) (int, error) {
	var userID int
	query := `
		SELECT user_id FROM refresh_token 
		WHERE token = $1 AND expire_at > NOW()
	`
	err := r.db.QueryRow(query, token).Scan(&userID)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (r *userRepository) DeleteRefreshToken(token string) error {
	_, err := r.db.Exec("DELETE FROM refresh_token WHERE token = $1", token)
	return err
}
