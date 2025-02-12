package usecases

import (
	"context"

	"github.com/codepnw/react_go_ecom/internal/entities"
	"github.com/codepnw/react_go_ecom/internal/repositories"
)

type CategoryUsecase interface {
	CreateCategory(ctx context.Context, title string) error
	ListCategory(ctx context.Context) ([]*entities.Category, error)
	DeleteCategory(ctx context.Context, id int) error
}

type categoryUsecase struct {
	repo repositories.CategoryRepo
}

func NewCategoryUsecase(repo repositories.CategoryRepo) CategoryUsecase {
	return &categoryUsecase{repo: repo}
}

func (uc *categoryUsecase) CreateCategory(ctx context.Context, title string) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	return uc.repo.Create(ctx, title)
}

func (uc *categoryUsecase) ListCategory(ctx context.Context) ([]*entities.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	return uc.repo.List(ctx)
}

func (uc *categoryUsecase) DeleteCategory(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	return uc.repo.Delete(ctx, id)
}
