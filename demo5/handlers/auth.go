package handlers

import (
	"time"

	"github.com/addetz/secure-code-go/demo4/data"
	"github.com/addetz/secure-code-go/demo4/db"
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

func NewUserAuthService(secret string, dbService db.DatabaseService) *UserAuthService {
	us := data.NewUserService(dbService)
	ns := data.NewSecretNoteService(dbService)
	return &UserAuthService{
		userService:        us,
		secretNotesService: ns,
		secret:             secret,
	}
}

func (us *UserAuthService) EncodeToken(username string) (string, error) {
	// Set custom claims
	claims := &JWTCustomClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	return token.SignedString([]byte(us.secret))
}
