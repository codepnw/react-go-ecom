package storage

import (
	"database/sql"

	"github.com/codepnw/react_go_ecom/internal/handlers"
	"github.com/codepnw/react_go_ecom/internal/repositories"
	"github.com/codepnw/react_go_ecom/internal/usecases"
)

type Storage struct {
	Category handlers.CategoryHandler
	Product  handlers.ProductHandler
}

func NewStorage(db *sql.DB) Storage {
	catRepo := repositories.NewCategoryRepo(db)
	catUc := usecases.NewCategoryUsecase(catRepo)
	catHandler := handlers.NewCategoryHandler(catUc)

	proRepo := repositories.NewProductRepository(db)
	proUsecase := usecases.NewProductUsecase(proRepo)
	proHandler := handlers.NewProductHandler(proUsecase)

	return Storage{
		Category: catHandler,
		Product:  proHandler,
	}
}
