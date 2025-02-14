package usecases

import (
	"context"

	"github.com/codepnw/react_go_ecom/internal/entities"
	"github.com/codepnw/react_go_ecom/internal/repositories"
)

type ProductUsecase interface {
	Create(ctx context.Context, req *entities.ProductPayloadReq) error
	GetByID(ctx context.Context, id int) (*entities.Product, error)
	List(ctx context.Context) ([]*entities.Product, error)
	Update(ctx context.Context, id int, req entities.Product) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, text string) ([]*entities.Product, error)
}

type productUsecase struct {
	repo repositories.ProductRepository
}

func NewProductUsecase(repo repositories.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}

func (uc *productUsecase) Create(ctx context.Context, req *entities.ProductPayloadReq) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	return uc.repo.Create(ctx, req)
}

func (uc *productUsecase) GetByID(ctx context.Context, id int) (*entities.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	product, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (uc *productUsecase) List(ctx context.Context) ([]*entities.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	products, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (uc *productUsecase) Update(ctx context.Context, id int, req entities.Product) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	if err := uc.repo.Update(ctx, id, req); err != nil {
		return err
	}

	return nil
}

func (uc *productUsecase) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	return uc.repo.Delete(ctx, id)
}

func (uc *productUsecase) Search(ctx context.Context, text string) ([]*entities.Product, error) {
	return uc.repo.Search(ctx, text)
}