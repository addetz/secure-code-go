package handlers

import (
	"github.com/addetz/secure-code-go/demo3/data"
	"github.com/golang-jwt/jwt/v5"
)

// UserRequest represents a login or sign up request
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWTCustomClaims are custom claims extending default ones.
type JWTCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UserAuthService struct {
	userService        *data.UserService
	secretNotesService *data.SecretNoteService
	secret             string
}

func NewUserAuthService(secret string) *UserAuthService {
	us := data.NewUserService()
	ns := data.NewSecretNoteService()
	return &UserAuthService{
		userService:        us,
		secretNotesService: ns,
		secret:             secret,
	}
}
