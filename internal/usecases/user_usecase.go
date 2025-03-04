package usecases

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/codepnw/react_go_ecom/config"
	"github.com/codepnw/react_go_ecom/internal/entities"
	"github.com/codepnw/react_go_ecom/internal/repositories"
	"github.com/codepnw/react_go_ecom/internal/utils"
	"github.com/codepnw/react_go_ecom/pkg/auth"
)

type UserUsecase interface {
	Register(ctx context.Context, req *entities.UserRegisterReq) (*entities.User, error)
	Login(ctx context.Context, req *entities.UserLoginReq) (string, string, error)
	GetProfile(ctx context.Context, id string) (*entities.User, error)
	RefreshToken(refreshToken string) (string, error)
	Logout(token string) error
}

type userUsecase struct {
	repo repositories.UserRepository
	cfg  config.JWTConfig
}

func NewUserUsecase(repo repositories.UserRepository, cfg config.JWTConfig) UserUsecase {
	return &userUsecase{
		repo: repo,
		cfg:  cfg,
	}
}

func (uc *userUsecase) Register(ctx context.Context, req *entities.UserRegisterReq) (*entities.User, error) {
	var user entities.User

	if !user.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email address")
	}

	if !user.ValidatePassword(req.Password) {
		return nil, errors.New("password least 6 character")
	}

	hashedPassword, err := user.HashedPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user = entities.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RoleID:    3,
		Enabled:   true,
		Address:   req.Address,
		CreatedAt: utils.ThaiTime,
	}

	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	userID, err := uc.repo.Create(ctx, &user)
	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"":
			return nil, errors.New("email is already exists")
		default:
			return nil, err
		}
	}

	u, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (uc *userUsecase) Login(ctx context.Context, req *entities.UserLoginReq) (string, string, error) {
	var u entities.User

	if !u.ValidateEmail(req.Email) {
		return "", "", errors.New("invalid email address")
	}

	if !u.ValidatePassword(req.Password) {
		return "", "", errors.New("password least 6 character")
	}

	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	user, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", errors.New("email or password is invalid")
		}
		return "", "", err
	}

	if err := user.CompareHashedPassword(user.Password, req.Password); err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := auth.GenerateToken(user.ID, uc.cfg)
	if err != nil {
		return "", "", err
	}

	// Save Refresh Token
	tokenExpire := uc.cfg.RefreshTokenExpire
	expireAt := time.Now().Add(time.Duration(tokenExpire) * time.Minute)
	if err := uc.repo.SaveRefreshToken(user.ID, refreshToken, expireAt); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (uc *userUsecase) RefreshToken(refreshToken string) (string, error) {
	// claims, err := auth.ValidateToken(refreshToken, uc.cfg.Secret)
	// if err != nil {
	// 	log.Println(err)
	// 	return "", errors.New("invalid refresh token")
	// }

	// Check Refresh Token
	userID, err := uc.repo.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid or expires refresh token")
	}

	// New Access Token
	idStr := strconv.Itoa(userID)
	newAccessToken, _, err := auth.GenerateToken(idStr, uc.cfg)
	if err != nil {
		log.Println(err)
		return "", errors.New("failed to generate token")
	}

	return newAccessToken, nil
}

func (uc *userUsecase) GetProfile(ctx context.Context, id string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, contextTimeoutQuery)
	defer cancel()

	log.Println("user_id:", id)

	return uc.repo.GetByID(ctx, id)
}

func (uc *userUsecase) Logout(token string) error {
	return uc.repo.DeleteRefreshToken(token)
}
