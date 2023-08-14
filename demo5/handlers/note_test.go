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
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNotes(t *testing.T) {
	username := "user1"
	noteText := "my super duper secret"
	mockDB := new(mocks.DatabaseServiceMock)
	userAuthService := handlers.NewUserAuthService("testing-signing-key", mockDB)
	mockDB.On("GetUser", username).Return(nil, nil)
	token, err := userAuthService.EncodeToken(username)
	assert.Nil(t, err)
	mockDB.On("AddNote", mock.AnythingOfType("string"), username, noteText).Return(nil)
	mockDB.On("GetUserNotes", username).Return([]db.Note{
		{
			ID:       uuid.New().String(),
			Username: username,
			Text:     noteText,
		},
	}, nil)

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
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Contains(t, res.Body.String(), username)
		assert.Contains(t, res.Body.String(), noteText)
	})

	t.Run("successful get notes", func(t *testing.T) {
		reqGet := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/secretNotes/%s", username), nil)
		reqGet.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		reqGet.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
		resGet := httptest.NewRecorder()
		e.ServeHTTP(resGet, reqGet)

		assert.Equal(t, http.StatusOK, resGet.Code)
		assert.Contains(t, resGet.Body.String(), username)
		assert.Contains(t, resGet.Body.String(), noteText)
	})

	t.Run("no token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/secretNotes/%s", username), nil)
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Code)
		assert.Contains(t, res.Body.String(), "missing or malformed jwt")
	})
}
