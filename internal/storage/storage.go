package storage

import (
	"database/sql"

	"github.com/codepnw/react_go_ecom/internal/handlers"
	"github.com/codepnw/react_go_ecom/internal/repositories"
	"github.com/codepnw/react_go_ecom/internal/usecases"
)

type Storage struct {
	Category handlers.CategoryHandler
}

func NewStorage(db *sql.DB) Storage {
	catRepo := repositories.NewCategoryRepo(db)
	catUc := usecases.NewCategoryUsecase(catRepo)
	catHandler := handlers.NewCategoryHandler(catUc)

	return Storage{
		Category: catHandler,
	}
}