package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/addetz/secure-code-go/demo4/db"
	"github.com/addetz/secure-code-go/demo4/handlers"
	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
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

	// Connect to database
	dbConn := connectDatabase()
	// Shut down connection when server shuts down
	defer func() {
		dbConn.Close()
	}()
	dbService := db.NewDatabaseService(dbConn)

	// Set up internal services
	userAuthService := handlers.NewUserAuthService(signingKey, dbService)

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
	// Get all user's notes
	r.GET("/secretNotes/:id", func(c echo.Context) error {
		return userAuthService.GetUserNotes(c)
	})
	// Add new note for user
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

func connectDatabase() *sql.DB {
	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		log.Fatal("POSTGRES_USER variable must be set")
	}
	pwd, ok := os.LookupEnv("POSTGRES_PWD")
	if !ok {
		log.Fatal("POSTGRES_PWD variable must be set")
	}
	db, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		log.Fatal("POSTGRES_DB variable must be set")
	}

	connectionStr := fmt.Sprintf("user=%s password=%s dbname=%s", user, pwd, db)
	conn, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal("connection error", err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal("ping error", err)
	}
	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS users (username VARCHAR(50) PRIMARY KEY, pwd VARCHAR(100) NOT NULL)")
	if err != nil {
		log.Fatal("create users", err)
	}
	_, err = conn.Exec("CREATE TABLE IF NOT EXISTS notes ( id VARCHAR (50) PRIMARY KEY," +
		"username VARCHAR(50) REFERENCES users (username), noteText VARCHAR (500) NOT NULL)")
	if err != nil {
		panic(err)
	}
	return conn
}
