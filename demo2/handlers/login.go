package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (authService *UserAuthService) Login(c echo.Context) error {
	u := new(UserRequest)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err := authService.userService.ValidatePassword(u.Username, u.Password); err != nil {
		return errors.Wrap(err, "login")
	}

	// Set custom claims
	claims := &JWTCustomClaims{
		u.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(authService.secret))
	if err != nil {
		return errors.Wrap(err, "login")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
