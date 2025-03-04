package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/codepnw/react_go_ecom/internal/entities"
	"github.com/codepnw/react_go_ecom/internal/usecases"
	"github.com/codepnw/react_go_ecom/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Profile(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}

type userHandler struct {
	uc usecases.UserUsecase
}

func NewUserHandler(uc usecases.UserUsecase) UserHandler {
	return &userHandler{uc: uc}
}

func (h *userHandler) Register(c *gin.Context) {
	var req entities.UserRegisterReq

	if err := c.ShouldBind(&req); err != nil {
		utils.NewResponse(c).Error(http.StatusBadRequest, err)
		return
	}

	user, err := h.uc.Register(c.Request.Context(), &req)
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusCreated, user)
}

func (h *userHandler) Login(c *gin.Context) {
	var req entities.UserLoginReq

	if err := c.ShouldBind(&req); err != nil {
		utils.NewResponse(c).Error(http.StatusBadRequest, err)
		return
	}

	accessToken, refreshToken, err := h.uc.Login(c.Request.Context(), &req)
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	// Set Header
	c.Header("Authorizarion", fmt.Sprintf("Bearer %s", accessToken))

	utils.NewResponse(c).Success(http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *userHandler) Profile(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		utils.NewResponse(c).Error(http.StatusBadRequest, errors.New("user_id not found"))
		return
	}

	user, err := h.uc.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	log.Println("Profile:", user)

	utils.NewResponse(c).Success(http.StatusOK, user)
}

func (h *userHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBind(&req); err != nil {
		utils.NewResponse(c).Error(http.StatusBadRequest, err)
		return
	}

	newAccessToken, err := h.uc.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, newAccessToken)
}

func (h *userHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBind(&req); err != nil {
		utils.NewResponse(c).Error(http.StatusBadRequest, err)
		return
	}

	if err := h.uc.Logout(req.RefreshToken); err != nil {
		utils.NewResponse(c).Error(http.StatusInternalServerError, err)
		return
	}

	utils.NewResponse(c).Success(http.StatusNoContent, "user logout!")
}
