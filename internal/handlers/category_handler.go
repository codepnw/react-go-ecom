package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codepnw/react_go_ecom/internal/entities"
	"github.com/codepnw/react_go_ecom/internal/usecases"
	"github.com/gin-gonic/gin"
)

type CategoryHandler interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	Delete(c *gin.Context)
}

type categoryHandler struct {
	uc usecases.CategoryUsecase
}

func NewCategoryHandler(uc usecases.CategoryUsecase) CategoryHandler {
	return &categoryHandler{uc: uc}
}

func (h *categoryHandler) Create(c *gin.Context) {
	var payload entities.CategoryReq

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.CreateCategory(c.Request.Context(), payload.Title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "added category"})
}

func (h *categoryHandler) List(c *gin.Context) {
	categories, err := h.uc.ListCategory(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *categoryHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.uc.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": fmt.Sprintf("category id %v deleted", id)})
}
