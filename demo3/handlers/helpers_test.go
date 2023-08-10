package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/addetz/secure-code-go/demo3/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type token struct {
	Token string `json:"token"`
}

func signUp(t *testing.T, userAuthService *handlers.UserAuthService) (*token, string) {
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

	return newToken, "user1"
}
