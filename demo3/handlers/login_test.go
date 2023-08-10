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

func TestLogin(t *testing.T) {
	successfulUser := `{"username":"user1","password":"potato-cheese-entropy-romania"}`
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(successfulUser))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	userAuthService := handlers.NewUserAuthService("testing-signing-key")

	// Assertions Setup
	err := userAuthService.SignUp(c)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "token")

	// Login
	reqLogin := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(successfulUser))
	reqLogin.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recLogin := httptest.NewRecorder()
	cLogin := e.NewContext(reqLogin, recLogin)
	err = userAuthService.Login(cLogin)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, recLogin.Code)
	assert.Contains(t, recLogin.Body.String(), "token")
}
