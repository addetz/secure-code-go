package main_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouterHTTP(t *testing.T) {
	// Send an HTTP request
	port := readPort()
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s", port))
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRouterHTTPS(t *testing.T) {
	// Send an HTTPS request
	port := readPort()
	resp, err := http.Get(fmt.Sprintf("https://localhost:%s", port))
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func readPort() string {
	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		return "1232"
	}
	return port
}
