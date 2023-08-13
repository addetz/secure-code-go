package handlers

import (
	"net/http"

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

	t, err := authService.EncodeToken(u.Username)
	if err != nil {
		return errors.Wrap(err, "login")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
