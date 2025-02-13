package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/codepnw/react_go_ecom/internal/entities"
	"github.com/codepnw/react_go_ecom/internal/utils"
)

type ProductRepository interface {
	Create(ctx context.Context, req *entities.ProductPayloadReq) error
	GetByID(ctx context.Context, id int) (*entities.Product, error)
	List(ctx context.Context) ([]*entities.Product, error)
	Update(ctx context.Context, id int, req entities.Product) error
	Delete(ctx context.Context, id int) error
}

type productRepository struct {
	db             *sql.DB
	queryFields    []string
	values         []any
	lastStackIndex int
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		db:          db,
		queryFields: make([]string, 0),
		values:      make([]any, 0),
	}
}

func (r *productRepository) Create(ctx context.Context, req *entities.ProductPayloadReq) error {
	query := `
		INSERT INTO products (title, description, price, sold, quantity)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int

	err := r.db.QueryRowContext(
		ctx,
		query,
		&req.Title,
		&req.Description,
		&req.Price,
		&req.Sold,
		&req.Quantity,
	).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*entities.Product, error) {
	query := `
		SELECT id, title, description, price, sold, quantity, created_at, updated_at
		FROM products WHERE id = $1
	`
	var p entities.Product

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Title,
		&p.Description,
		&p.Price,
		&p.Sold,
		&p.Quantity,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) List(ctx context.Context) ([]*entities.Product, error) {
	query := `
		SELECT id, title, description, price, sold, quantity, created_at, updated_at
		FROM products
	`
	var products []*entities.Product

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p entities.Product
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.Price,
			&p.Sold,
			&p.Quantity,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) Update(ctx context.Context, id int, req entities.Product) error {
	var fields []string
	var values []any
	var lastIndex int

	if req.Title != "" {
		values = append(values, req.Title)
		lastIndex = len(values)

		fields = append(fields, fmt.Sprintf("title = $%d", lastIndex))
	}

	if req.Description != "" {
		values = append(values, req.Description)
		lastIndex = len(values)

		fields = append(fields, fmt.Sprintf("description = $%d", lastIndex))
	}

	if req.Price != 0 {
		values = append(values, req.Price)
		lastIndex = len(values)

		fields = append(fields, fmt.Sprintf("price = $%d", lastIndex))
	}

	if req.Sold != 0 {
		values = append(values, req.Sold)
		lastIndex = len(values)

		fields = append(fields, fmt.Sprintf("sold = $%d", lastIndex))
	}

	if req.Quantity != 0 {
		values = append(values, req.Quantity)
		lastIndex = len(values)

		fields = append(fields, fmt.Sprintf("quantity = $%d", lastIndex))
	}

	// Add Field updated_at
	values = append(values, utils.ThaiTime)
	lastIndex = len(values)
	fields = append(fields, fmt.Sprintf("updated_at = $%d", lastIndex))

	// Add Product ID
	values = append(values, id)
	lastIndex = len(values)
	query := fmt.Sprintf("UPDATE products SET %s WHERE id = $%d", strings.Join(fields, ", "), lastIndex)

	log.Println(query)

	result, err := r.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}
