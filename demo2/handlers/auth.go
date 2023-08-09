package handlers

import (
	"github.com/addetz/secure-code-go/demo2/data"
	"github.com/golang-jwt/jwt/v5"
)

// UserRequest represents a login or sign up request
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWTCustomClaims are custom claims extending default ones.
type JWTCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type UserAuthService struct {
	userService *data.UserService
	secret      string
}

func NewUserAuthService(secret string) *UserAuthService {
	us := data.NewUserService()
	return &UserAuthService{
		userService: us,
		secret:      secret,
	}
}
