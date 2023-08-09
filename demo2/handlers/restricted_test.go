package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/addetz/secure-code-go/demo2/handlers"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type token struct {
	Token string `json:"token"`
}

func TestRestricted(t *testing.T) {
	userAuthService := handlers.NewUserAuthService("testing-signing-key")
	token := signUp(t, userAuthService)

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
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token.Token))
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

func signUp(t *testing.T, userAuthService *handlers.UserAuthService) *token {
	t.Helper()
	successfulUser := `{"username":"user1","password":"potato-cheese-entropy-romania"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(successfulUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	err := userAuthService.SignUp(c)
	assert.Nil(t, err)
	newToken := new(token)
	assert.Nil(t, json.Unmarshal(rec.Body.Bytes(), newToken))

	return newToken
}
