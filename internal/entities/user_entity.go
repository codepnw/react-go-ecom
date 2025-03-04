package entities

import (
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	RoleID    int        `json:"role_id"`
	Enabled   bool       `json:"enabled"`
	Address   string     `json:"address"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type UserRegisterReq struct {
	Email     string `form:"email" json:"email" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	FirstName string `form:"first_name" json:"first_name" binding:"required"`
	LastName  string `form:"last_name" json:"last_name" binding:"required"`
	RoleID    int    `form:"role_id" json:"role_id"`
	Enabled   bool   `form:"enabled" json:"enabled"`
	Address   string `form:"address" json:"address" binding:"required"`
}

type UserLoginReq struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (u *User) HashedPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (u *User) CompareHashedPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) ValidateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

func (u *User) ValidatePassword(password string) bool {
	return len(password) >= 6
}
