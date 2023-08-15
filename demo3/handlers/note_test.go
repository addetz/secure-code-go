package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/addetz/secure-code-go/demo3/handlers"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNotes(t *testing.T) {
	userAuthService := handlers.NewUserAuthService("testing-signing-key")
	token, username := signUp(t, userAuthService)

	// set up restricted path middleware
	e := echo.New()
	e.POST("/secretNotes/:id", func(c echo.Context) error {
		return userAuthService.AddUserNote(c)
	})
	e.GET("/secretNotes/:id", func(c echo.Context) error {
		return userAuthService.GetUserNotes(c)
	})

	e.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.JWTCustomClaims)
		},
		SigningKey: []byte("testing-signing-key"),
	}))

	t.Run("successful add note", func(t *testing.T) {
		newNote := `{"text":"my super duper secret"}`
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/secretNotes/%s", username), strings.NewReader(newNote))
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token.Token))
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), username)
		assert.Contains(t, res.Body.String(), "my super duper secret")
	})

	t.Run("successful get notes", func(t *testing.T) {
		newNote := `{"text":"my super duper secret"}`
		reqPost := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/secretNotes/%s", username), strings.NewReader(newNote))
		reqPost.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token.Token))
		resPost := httptest.NewRecorder()
		e.ServeHTTP(resPost, reqPost)

		assert.Equal(t, http.StatusOK, resPost.Code)
		assert.Contains(t, resPost.Body.String(), username)
		assert.Contains(t, resPost.Body.String(), "my super duper secret")

		reqGet := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/secretNotes/%s", username), nil)
		reqGet.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token.Token))
		resGet := httptest.NewRecorder()
		e.ServeHTTP(resGet, reqGet)

		assert.Equal(t, http.StatusOK, resGet.Code)
		assert.Contains(t, resGet.Body.String(), username)
		assert.Contains(t, resGet.Body.String(), "my super duper secret")
	})

	t.Run("no token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/secretNotes/%s", username), nil)
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Contains(t, res.Body.String(), "missing or malformed jwt")
	})
}
