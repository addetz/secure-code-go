package handlers

import (
	"net/http"

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

	t, err := authService.EncodeToken(u.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"token": t,
	})
}
