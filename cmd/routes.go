package main

import (
	"database/sql"
	"net/http"

	"github.com/codepnw/react_go_ecom/internal/storage"
	"github.com/gin-gonic/gin"
)

func apiRoutes(db *sql.DB, r *gin.Engine, version string) {
	store := storage.NewStorage(db)

	router := r.Group("/api/" + version)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API running !"})
	})

	catRouter := router.Group("/categories")

	catRouter.POST("/", store.Category.Create)
	catRouter.GET("/", store.Category.List)
	catRouter.DELETE("/:id", store.Category.Delete)

}
