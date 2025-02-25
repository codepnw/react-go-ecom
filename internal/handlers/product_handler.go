package handlers

import (
	"fmt"
	"net/http"

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
	ProductPurchase(c *gin.Context)
	CheckOutOfStock(c *gin.Context)
	RestockProduct(c *gin.Context)
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

	id, err := h.uc.Create(c.Request.Context(), &req)
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	product, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusCreated, product)
}

func (h *productHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, product)
}

func (h *productHandler) List(c *gin.Context) {
	search := c.Query("search")
	limit := c.Query("limit")
	offset := c.Query("offset")

	if search != "" {
		products, err := h.uc.Search(c.Request.Context(), search)
		if err != nil {
			utils.NewResponse(c).Error(http.StatusInternalServerError, err)
			return
		}

		utils.NewResponse(c).Success(http.StatusOK, products)
		return
	}

	products, err := h.uc.List(c.Request.Context(), limit, offset)
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, products)
}

func (h *productHandler) Update(c *gin.Context) {
	id := c.Param("id")

	req := entities.Product{}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewResponse(c).Error(http.StatusBadRequest, err)
		return
	}

	if err := h.uc.Update(c.Request.Context(), id, req); err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	product, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, product)
}

func (h *productHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, fmt.Sprintf("product_id %s deleted", id))
}

func (h *productHandler) ProductPurchase(c *gin.Context) {
	req := entities.ProductStock{}

	if err := c.ShouldBind(&req); err != nil {
		utils.NewResponse(c).Error(http.StatusBadRequest, err)
		return
	}

	if err := h.uc.PurchaseProduct(&req); err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, "purchase successful")
}

func (h *productHandler) CheckOutOfStock(c *gin.Context) {
	products, err := h.uc.CheckOutOfStock()
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, products)
}

func (h *productHandler) RestockProduct(c *gin.Context) {
	req := entities.ProductStock{}

	if err := c.ShouldBind(&req); err != nil {
		utils.NewResponse(c).Error(http.StatusBadRequest, err)
		return
	}

	if err := h.uc.RestockProduct(&req); err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, fmt.Sprintf("product_id %s added %d stock", req.ProductID, req.Quantity))
}