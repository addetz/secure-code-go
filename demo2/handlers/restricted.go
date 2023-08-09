package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (authService *UserAuthService) RestrictedPath(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTCustomClaims)
	name := claims.Name
	if err := authService.userService.ValidateUser(name); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": fmt.Sprintf("You're logged in %s!", name),
	})
}
