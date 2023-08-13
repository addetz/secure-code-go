package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/addetz/secure-code-go/demo3/handlers"
	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/golang-jwt/jwt/v5"
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
	signingKey, ok := os.LookupEnv("SIGNING_KEY")
	if !ok {
		log.Fatal("SIGNING_KEY variable must be set")
	}
	// Read port if one is set
	port := readPort()

	// Set up internal services
	userAuthService := handlers.NewUserAuthService(signingKey)

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
			Message: "Hello, Gophers!",
		})
	})

	// Set up authentication routes
	e.POST("/signup", func(c echo.Context) error {
		return userAuthService.SignUp(c)
	})
	e.POST("/login", func(c echo.Context) error {
		return userAuthService.Login(c)
	})

	// Restricted route only for logged in users
	r := e.Group("/restricted")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.JWTCustomClaims)
		},
		SigningKey: []byte(signingKey),
	}

	r.Use(echojwt.WithConfig(config))
	r.GET("", func(c echo.Context) error {
		return userAuthService.RestrictedPath(c)
	})
	r.GET("/secretNotes/:id", func(c echo.Context) error {
		return userAuthService.GetUserNotes(c)
	})
	r.POST("/secretNotes/:id", func(c echo.Context) error {
		return userAuthService.AddUserNote(c)
	})

	log.Printf("Listening on :%s...\n", port)
	if err := s.ListenAndServeTLS(certFile, keyFile); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func readPort() string {
	port, ok := os.LookupEnv("SERVER_PORT")
	if !ok {
		return "1323"
	}
	return port
}
