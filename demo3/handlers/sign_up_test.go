package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/addetz/secure-code-go/demo3/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)


func TestSignUp(t *testing.T) {
	successfulUser := `{"username":"user1","password":"potato-cheese-entropy-romania"}`

	t.Run("successful sign up", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(successfulUser))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		userAuthService := handlers.NewUserAuthService("testing-signing-key")

		// Assertions
		err := userAuthService.SignUp(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")
	})

	t.Run("repeated sign up", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(successfulUser))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		userAuthService := handlers.NewUserAuthService("testing-signing-key")

		// Assertions
		err := userAuthService.SignUp(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")

		err = userAuthService.SignUp(c)
		assert.NotNil(t, err)
	})
}
