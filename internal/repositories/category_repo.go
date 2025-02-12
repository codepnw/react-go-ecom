package repositories

import (
	"context"
	"database/sql"

	"github.com/codepnw/react_go_ecom/internal/entities"
)

type CategoryRepo interface {
	Create(ctx context.Context, title string) error
	List(ctx context.Context) ([]*entities.Category, error)
	Delete(ctx context.Context, id int) error
}

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) CategoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) Create(ctx context.Context, title string) error {
	query := `INSERT INTO categories (title) VALUES ($1)`

	_, err := r.db.ExecContext(ctx, query, title)
	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepo) List(ctx context.Context) ([]*entities.Category, error) {
	query := `SELECT id, title, created_at FROM categories`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*entities.Category
	for rows.Next() {
		var cat entities.Category
		if err := rows.Scan(&cat.ID, &cat.Title, &cat.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, &cat)
	}

	return categories, nil
}

func (r *categoryRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}
