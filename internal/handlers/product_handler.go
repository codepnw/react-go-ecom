package handlers

import (
	"net/http"
	"strconv"

	"github.com/codepnw/react_go_ecom/internal/entities"
	"github.com/codepnw/react_go_ecom/internal/usecases"
	"github.com/codepnw/react_go_ecom/internal/utils"
	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type productHandler struct {
	uc usecases.ProductUsecase
}

func NewProductHandler(uc usecases.ProductUsecase) ProductHandler {
	return &productHandler{uc: uc}
}

func (h *productHandler) Create(c *gin.Context) {
	var req entities.ProductPayloadReq

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewResponse(c).Error(http.StatusBadRequest, err)
		return
	}

	if err := h.uc.Create(c.Request.Context(), &req); err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusCreated, "product created")
}

func (h *productHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	product, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, product)
}

func (h *productHandler) List(c *gin.Context) {
	search := c.Query("search")

	if search != "" {
		products, err := h.uc.Search(c.Request.Context(), search)
		if err != nil {
			utils.NewResponse(c).Error(http.StatusInternalServerError, err)
			return
		}

		utils.NewResponse(c).Success(http.StatusOK, products)
		return
	}

	products, err := h.uc.List(c.Request.Context())
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, products)
}

func (h *productHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	req := entities.Product{}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewResponse(c).Error(http.StatusBadRequest, err)
		return
	}

	if err := h.uc.Update(c.Request.Context(), id, req); err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, "product updated")
}

func (h *productHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusNoContent, "product deleted")
}
