package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func (authService *UserAuthService) SignUp(c echo.Context) error {
	// Read user request
	u := new(UserRequest)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	// Send user data to the user service
	if err := authService.userService.Add(u.Username, u.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "sign up"))
	}

	// Set custom claims
	claims := &JWTCustomClaims{
		u.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	// Create token with claims
	// See a token:
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkZGV0eiIsImV4cCI6MTY5MTk2NzA5MH0.QFvlhQ7oSk6gSwfa_AOxUyAxD10pqQfXhOUtf4YLu5w
	// Decoded:
	// {"alg":"HS256","typ":"JWT"}{"username":"addetz","exp":1691967090}aC�TCJjA9K_
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(authService.secret))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"token": t,
	})
}
