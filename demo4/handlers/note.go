package handlers

import (
	"errors"
	"net/http"

	"github.com/addetz/secure-code-go/demo4/data"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetUserNotes returns all the notes of a given user.
func (authService *UserAuthService) GetUserNotes(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTCustomClaims)
	name := claims.Username
	if err := authService.userService.ValidateUser(name); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	paramName := c.Param("id")
	if name != paramName {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("not logged in as notes owner"))
	}
	secretNotes, err := authService.secretNotesService.GetAll(paramName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"username": name,
		"notes":    secretNotes,
	})
}

// AddUserNote adds a note belonging to the given user
func (authService *UserAuthService) AddUserNote(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTCustomClaims)
	name := claims.Username
	if err := authService.userService.ValidateUser(name); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}
	paramName := c.Param("id")
	if name != paramName {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("not logged in as notes owner"))
	}

	newNote := new(data.SecretNote)
	if err := c.Bind(newNote); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	//add the note
	authService.secretNotesService.Add(paramName, *newNote)
	secretNotes, err := authService.secretNotesService.GetAll(paramName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"username": name,
		"notes":    secretNotes,
	})
}
