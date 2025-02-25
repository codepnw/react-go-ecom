package usecases

import (
	"context"

	"github.com/codepnw/react_go_ecom/internal/entities"
	"github.com/codepnw/react_go_ecom/internal/repositories"
	"github.com/codepnw/react_go_ecom/internal/utils"
)

type ProductUsecase interface {
	Create(ctx context.Context, req *entities.ProductPayloadReq) (string, error)
	GetByID(ctx context.Context, id string) (*entities.Product, error)
	List(ctx context.Context, limit, offset string) ([]*entities.Product, error)
	Update(ctx context.Context, id string, req entities.Product) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, text string) ([]*entities.Product, error)
}

type productUsecase struct {
	repo repositories.ProductRepository
}

func NewProductUsecase(repo repositories.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}

func (uc *productUsecase) Create(ctx context.Context, req *entities.ProductPayloadReq) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	var product = &entities.Product{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   utils.ThaiTime,
	}

	return uc.repo.Create(ctx, product)
}

func (uc *productUsecase) GetByID(ctx context.Context, id string) (*entities.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	return uc.repo.GetByID(ctx, id)
}

func (uc *productUsecase) List(ctx context.Context, limit, offset string) ([]*entities.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	if limit == "" {
		limit = "10"
	}

	if offset == "" {
		offset = "0"
	}

	products, err := uc.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (uc *productUsecase) Update(ctx context.Context, id string, req entities.Product) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	if err := uc.repo.Update(ctx, id, req); err != nil {
		return err
	}

	return nil
}

func (uc *productUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	return uc.repo.Delete(ctx, id)
}

func (uc *productUsecase) Search(ctx context.Context, text string) ([]*entities.Product, error) {
	return uc.repo.Search(ctx, text)
}
