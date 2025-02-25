package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/codepnw/react_go_ecom/internal/entities"
	"github.com/codepnw/react_go_ecom/internal/utils"
)

type ProductRepository interface {
	Create(ctx context.Context, req *entities.Product) (string, error)
	GetByID(ctx context.Context, id string) (*entities.Product, error)
	List(ctx context.Context, limit, offset string) ([]*entities.Product, error)
	Update(ctx context.Context, id string, req entities.Product) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, text string) ([]*entities.Product, error)
	ReduceStock(productID string, stock int) error
	CheckStock(productID string) (int, error)
	AddSoldQuantity(productID string, quantity int) error
	CheckOutOfStock() ([]*entities.Product, error)
	RestockProduct(req *entities.ProductStock) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, req *entities.Product) (string, error) {
	query := `
		INSERT INTO products (title, description, price, stock, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING product_id
	`
	var id string

	err := r.db.QueryRowContext(
		ctx,
		query,
		&req.Title,
		&req.Description,
		&req.Price,
		&req.Stock,
		&req.CreatedAt,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*entities.Product, error) {
	query := `
		SELECT product_id, title, description, price, stock, created_at, updated_at
		FROM products WHERE product_id = $1
	`
	var p entities.Product

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.Title,
		&p.Description,
		&p.Price,
		&p.Stock,
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

func (r *productRepository) List(ctx context.Context, limit, offset string) ([]*entities.Product, error) {
	query := `
		SELECT product_id, title, description, price, stock, quantity, created_at, updated_at
		FROM products
	`
	query += fmt.Sprintf("LIMIT %s OFFSET %s", limit, offset)

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
			&p.Stock,
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

func (r *productRepository) Update(ctx context.Context, id string, req entities.Product) error {
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

	if req.Stock != 0 {
		values = append(values, req.Stock)
		lastIndex = len(values)

		fields = append(fields, fmt.Sprintf("stock = $%d", lastIndex))
	}

	// Add Field updated_at
	values = append(values, utils.ThaiTime)
	lastIndex = len(values)
	fields = append(fields, fmt.Sprintf("updated_at = $%d", lastIndex))

	// Add Product ID
	values = append(values, id)
	lastIndex = len(values)
	query := fmt.Sprintf("UPDATE products SET %s WHERE product_id = $%d", strings.Join(fields, ", "), lastIndex)

	// log.Println(query)

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

func (r *productRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE product_id = $1`

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

func (r *productRepository) Search(ctx context.Context, text string) ([]*entities.Product, error) {
	query := `
		SELECT product_id, title, description, price, stock, quantity, created_at, updated_at
		FROM products WHERE title ILIKE $1 OR description ILIKE $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, "%"+text+"%", "%"+text+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entities.Product
	for rows.Next() {
		var p entities.Product
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.Quantity,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}

func (r *productRepository) ReduceStock(productID string, stock int) error {
	query := `
		UPDATE products SET stock = stock - $1
		WHERE product_id = $2 AND stock >= $1
	`
	result, err := r.db.Exec(query, stock, productID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("not enough stock")
	}

	return nil
}

func (r *productRepository) CheckStock(productID string) (int, error) {
	var stock int

	err := r.db.QueryRow("SELECT stock FROM products WHERE product_id = $1", productID).Scan(&stock)
	if err != nil {
		return 0, err
	}

	return stock, nil
}

func (r *productRepository) AddSoldQuantity(productID string, quantity int) error {
	result, err := r.db.Exec("UPDATE products SET quantity = $1 WHERE product_id = $2", quantity, productID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("update sold quantity fail")
	}

	return nil
}

func (r *productRepository) CheckOutOfStock() ([]*entities.Product, error) {
	query := `
		SELECT product_id, title, description, price, stock, quantity, created_at, updated_at
		FROM products WHERE stock = 0	
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*entities.Product{}

	for rows.Next() {
		p := entities.Product{}

		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.Quantity,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}

func (r *productRepository) RestockProduct(req *entities.ProductStock) error {
	query := `UPDATE products SET stock = stock + $1 WHERE product_id = $2`
	_, err := r.db.Exec(query, req.Quantity, req.ProductID)

	return err
}