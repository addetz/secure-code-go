package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/addetz/secure-code-go/demo4/handlers"
	"github.com/addetz/secure-code-go/demo4/mocks"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRestricted(t *testing.T) {
	user := "user1"
	mockDB := new(mocks.DatabaseServiceMock)
	userAuthService := handlers.NewUserAuthService("testing-signing-key", mockDB)
	mockDB.On("GetUser", user).Return(nil, nil)
	token, err := userAuthService.EncodeToken(user)
	assert.Nil(t, err)

	// set up restricted path middleware
	e := echo.New()
	e.GET("/restricted", func(c echo.Context) error {
		return userAuthService.RestrictedPath(c)
	})

	e.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.JWTCustomClaims)
		},
		SigningKey: []byte("testing-signing-key"),
	}))

	t.Run("successful restricted", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/restricted", nil)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		assert.Contains(t, res.Body.String(), "You're logged in")
	})

	t.Run("no token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/restricted", nil)
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Contains(t, res.Body.String(), "missing or malformed jwt")
	})
}
