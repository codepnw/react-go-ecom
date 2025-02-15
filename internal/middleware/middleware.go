package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/codepnw/react_go_ecom/config"
	"github.com/codepnw/react_go_ecom/internal/utils"
	"github.com/codepnw/react_go_ecom/pkg/auth"
	"github.com/gin-gonic/gin"
)

type middleware struct {
	cfg config.JWTConfig
}

func InitMiddleware(cfg config.JWTConfig) *middleware {
	return &middleware{cfg: cfg}
}

func (m *middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			utils.NewResponse(c).Error(http.StatusUnauthorized, errors.New("missing token"))
			c.Abort()
			return
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.NewResponse(c).Error(http.StatusUnauthorized, errors.New("invalid token format"))
			c.Abort()
			return
		}

		claims, err := auth.ValidateToken(parts[1], m.cfg.Secret)
		if err != nil {
			utils.NewResponse(c).Error(http.StatusUnauthorized, errors.New("invalid token"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.ID)
		c.Next()
	}
}
