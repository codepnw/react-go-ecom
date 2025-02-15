package main

import (
	"database/sql"
	"net/http"

	"github.com/codepnw/react_go_ecom/config"
	"github.com/codepnw/react_go_ecom/internal/middleware"
	"github.com/codepnw/react_go_ecom/internal/storage"
	"github.com/gin-gonic/gin"
)

func apiRoutes(db *sql.DB, cfg config.Config) *gin.Engine {
	r := gin.Default()

	store := storage.NewStorage(db, cfg)

	router := r.Group("/api/" + cfg.AppVersion)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API running !"})
	})

	// Auth Routes
	router.POST("/auth/register", store.User.Register)
	router.POST("/auth/login", store.User.Login)

	// Middleware Routes
	m := middleware.InitMiddleware(*cfg.JWTConfig)
	midRouter := router.Use(m.AuthMiddleware())
	midRouter.GET("/users/profile", store.User.Profile)

	// Categories Routes
	catRouter := router.Group("/categories")
	catRouter.POST("/", store.Category.Create)
	catRouter.GET("/", store.Category.List)
	catRouter.DELETE("/:id", store.Category.Delete)

	// Products Routes
	proRouter := router.Group("/products")
	proRouter.POST("/", store.Product.Create)
	proRouter.GET("/", store.Product.List)
	proRouter.GET("/:id", store.Product.GetByID)
	proRouter.PATCH("/:id", store.Product.Update)
	proRouter.DELETE("/:id", store.Product.Delete)

	return r
}
