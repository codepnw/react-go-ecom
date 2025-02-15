package auth

import (
	"errors"
	"time"

	"github.com/codepnw/react_go_ecom/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, cfg config.JWTConfig) (accessToken, refreshToken string, err error) {
	accessToken, err = generateToken(userID, []byte(cfg.Secret), time.Duration(cfg.AccessTokenExpire)*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = generateToken(userID, []byte(cfg.Secret), time.Duration(cfg.RefreshTokenExpire)*time.Minute)
	if err != nil {
		return "", "", err
	}

	return
}

func generateToken(userID string, secret []byte, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
