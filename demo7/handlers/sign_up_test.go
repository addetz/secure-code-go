package handlers_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/addetz/secure-code-go/demo4/handlers"
	"github.com/addetz/secure-code-go/demo4/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUp(t *testing.T) {
	password := "potato-cheese-entropy-romania"
	username := "user1"
	successfulUser := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)

	t.Run("successful sign up", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(successfulUser))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		mockDB := new(mocks.DatabaseServiceMock)
		userAuthService := handlers.NewUserAuthService("testing-signing-key", mockDB)
		mockDB.On("GetUser", username).Return(nil, errors.New("no user"))
		mockDB.On("AddUser", username, mock.AnythingOfType("string")).Return(nil)

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
		mockDB := new(mocks.DatabaseServiceMock)
		userAuthService := handlers.NewUserAuthService("testing-signing-key", mockDB)
		mockDB.On("GetUser", username).Return(nil, nil)

		// Assertions
		err := userAuthService.SignUp(c)
		assert.NotNil(t, err)
	})
}
