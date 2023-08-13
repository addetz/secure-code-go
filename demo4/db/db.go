package db

import (
	"database/sql"
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
	stmt, err := ds.db.Prepare("INSERT INTO users(username, pwd) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(username, pwd); err != nil {
		return err
	}
	return nil
}

// GetUser returns a user from the database or an error if none exists.
func (ds *dbService) GetUser(username string) (*User, error) {
	user := new(User)
	stmt, err := ds.db.Prepare("SELECT * FROM users WHERE username = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.QueryRow(username).Scan(user.Username, user.Pwd); err != nil {
		return nil, err
	}
	return user, nil
}

// AddNote creates a new note in the DB
func (ds *dbService) AddNote(id, username, text string) error {
	stmt, err := ds.db.Prepare("INSERT INTO notes(id, username, noteText) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err := stmt.Exec(id, username, text); err != nil {
		return err
	}
	return nil
}

// GetUserNotes returns all the notes of a given user from the database or an error.
func (ds *dbService) GetUserNotes(username string) ([]Note, error) {
	var notes []Note
	stmt, err := ds.db.Prepare("SELECT * FROM notes WHERE username = ?")
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
