package storage

import (
	"database/sql"

	"github.com/codepnw/react_go_ecom/config"
	"github.com/codepnw/react_go_ecom/internal/handlers"
	"github.com/codepnw/react_go_ecom/internal/repositories"
	"github.com/codepnw/react_go_ecom/internal/usecases"
)

type Storage struct {
	User     handlers.UserHandler
	Category handlers.CategoryHandler
	Product  handlers.ProductHandler
}

func NewStorage(db *sql.DB, cfg config.Config) Storage {
	userRepo := repositories.NewUserRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo, *cfg.JWTConfig)
	userHandler := handlers.NewUserHandler(userUsecase)

	catRepo := repositories.NewCategoryRepo(db)
	catUc := usecases.NewCategoryUsecase(catRepo)
	catHandler := handlers.NewCategoryHandler(catUc)

	proRepo := repositories.NewProductRepository(db)
	proUsecase := usecases.NewProductUsecase(proRepo)
	proHandler := handlers.NewProductHandler(proUsecase)

	return Storage{
		User:     userHandler,
		Category: catHandler,
		Product:  proHandler,
	}
}
