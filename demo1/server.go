package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const TIMEOUT = 3 * time.Second

type Response struct {
	Message string `json:"message"`
}

func main() {
	// Read paths to certificate & private key from environment variables
	certFile, ok := os.LookupEnv("SERVER_CERT_FILE")
	if !ok {
		log.Fatal("SERVER_CERT_FILE variable must be set")
	}
	keyFile, ok := os.LookupEnv("SERVER_KEY_FILE")
	if !ok {
		log.Fatal("SERVER_KEY_FILE variable must be set")
	}

	// Read port if one is set
	port := readPort()

	// Initialise echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configure server
	s := http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           e,
		ReadTimeout:       TIMEOUT,
		ReadHeaderTimeout: TIMEOUT,
		WriteTimeout:      TIMEOUT,
		IdleTimeout:       TIMEOUT,
	}

	// Set up the root route
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &Response{
			Message: "Hello, Gophers!\nGlad to see You there!",
		})
	})

	log.Printf("Listening on :%s...\n", port)
	if err := s.ListenAndServeTLS(certFile, keyFile); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func readPort() string {
	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		return "1232"
	}
	return port
}
