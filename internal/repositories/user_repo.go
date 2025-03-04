package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/codepnw/react_go_ecom/internal/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) (string, error)
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	SaveRefreshToken(userID string, token string, expires time.Time) error
	ValidateRefreshToken(token string) (int, error)
	DeleteRefreshToken(token string) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) (string, error) {
	query := `
		INSERT INTO users (email, password, first_name, last_name, role_id, enabled, images, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING user_id
	`
	var id string

	err := r.db.QueryRowContext(
		ctx,
		query,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.RoleID,
		&user.Enabled,
		"later",
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entities.User, error) {
	query := `
		SELECT user_id, email, first_name, last_name, role_id, address, enabled, created_at, updated_at
		FROM users WHERE user_id = $1
	`
	var user entities.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.RoleID,
		&user.Address,
		&user.Enabled,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT user_id, email, password, role_id
		FROM users WHERE email = $1
	`
	var user entities.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.RoleID,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) SaveRefreshToken(userID string, token string, expires time.Time) error {
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
