package db

import (
	"database/sql"
	"fmt"
	"log"
	"unicode"
)

type User struct {
	Username, Pwd string
}

type Note struct {
	ID, Username, Text string
}

type dbService struct {
	db *sql.DB
}

type DatabaseService interface {
	AddUser(username, pwd string) error
	GetUser(username string) (*User, error)
	AddNote(id, username, text string) error
	GetUserNotes(username string) ([]Note, error)
}

// NewDatabaseService initialises a DatabaseService given its dependencies.
func NewDatabaseService(db *sql.DB) *dbService {
	return &dbService{
		db: db,
	}
}

// AddUser creates a new user in the DB
func (ds *dbService) AddUser(username, pwd string) error {
	// let's assume, the username can only contain lowercase, ASCII chars, no numbers and be shorter than 10 chars
	err := validate(username)
	if err != nil {
		return err
	}

	tx, err := ds.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if _, err = tx.Exec("INSERT INTO users (username, pwd) VALUES(?, ?)", username, pwd); err != nil {
		return err
	}

	return nil
}

// GetUser returns a user from the database or an error if none exists.
func (ds *dbService) GetUser(username string) (*User, error) {
	var user User
	stmt, err := ds.db.Prepare("SELECT * FROM users WHERE username = $1 ")
	if err != nil {
		log.Println("error3", err)
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.QueryRow(username).Scan(&user.Username, &user.Pwd); err != nil {
		log.Println("error4", err)
		return nil, err
	}
	return &user, nil
}

// AddNote creates a new note in the DB
func (ds *dbService) AddNote(id, username, text string) error {
	stmt, err := ds.db.Prepare("INSERT INTO notes(id, username, noteText) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(id, username, text); err != nil {
		return err
	}
	return nil
}

func validate(str string) error {
	for _, r := range str {
		if unicode.IsUpper(r) {
			return fmt.Errorf("contains uppercase: %v", string(r))
		}

		if unicode.IsNumber(r) {
			return fmt.Errorf("contains number: %v", string(r))
		}

		if r > unicode.MaxASCII {
			return fmt.Errorf("contains non-ASCII char: %v", r)
		}
	}

	return nil
}

// GetUserNotes returns all the notes of a given user from the database or an error.
func (ds *dbService) GetUserNotes(username string) ([]Note, error) {
	var notes []Note
	stmt, err := ds.db.Prepare("SELECT * FROM notes WHERE username = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		n := Note{}
		if err := rows.Scan(&n.ID, &n.Username, &n.Text); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}
