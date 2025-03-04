package middleware

import (
	"database/sql"
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

func (m *middleware) RBACMiddleware(db *sql.DB, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.NewResponse(c).Error(http.StatusUnauthorized, errors.New("unauthorized"))
			c.Abort()
			return
		}

		var roleID int
		if err := db.QueryRow("SELECT role_id FROM users WHERE id = $1", userID); err != nil {
			utils.NewResponse(c).Error(http.StatusInternalServerError, errors.New("error role_id"))
			c.Abort()
			return
		}

		// Check Role Permission 
		query := `
			SELECT EXISTS (
				SELECT 1 FROM role_permissions
				JOIN permissions ON role_permission.permission_id = permission.id
				WHERE role_permission.role_id = $1 AND permission.permission_id = $2				
			)`
		var hasPermission bool
		err := db.QueryRow(query, roleID, permission).Scan(&hasPermission)
		if err != nil || !hasPermission {
			utils.NewResponse(c).Error(http.StatusForbidden, errors.New("forbidden"))
			c.Abort()
			return
		}

		c.Next()
	}
} 