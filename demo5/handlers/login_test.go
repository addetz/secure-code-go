package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/addetz/secure-code-go/demo4/db"
	"github.com/addetz/secure-code-go/demo4/handlers"
	"github.com/addetz/secure-code-go/demo4/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	password := "potato-cheese-entropy-romania"
	username := "user1"
	successfulUser := fmt.Sprintf(`{"username":"%s","password":"%s"}`,username, password) 
	expected, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Setup
	e := echo.New()
	mockDB := new(mocks.DatabaseServiceMock)
	userAuthService := handlers.NewUserAuthService("testing-signing-key", mockDB)
	mockDB.On("GetUser", username).Return(&db.User{
		Username: username,
		Pwd: string(expected),
	}, nil)

	// Login
	reqLogin := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(successfulUser))
	reqLogin.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recLogin := httptest.NewRecorder()
	cLogin := e.NewContext(reqLogin, recLogin)
	err := userAuthService.Login(cLogin)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, recLogin.Code)
	assert.Contains(t, recLogin.Body.String(), "token")
}
