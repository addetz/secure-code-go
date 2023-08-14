package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func FuzzAddNote(f *testing.F) {
	f.Add(10, "john", "🧞")

	db, mock, err := sqlmock.New()
	if err != nil {
		f.Fatal(err)
	}
	defer db.Close()

	// todo: finish fuzzing
	// refer to https://github.com/DATA-DOG/go-sqlmock
	// refer to https://go.dev/security/fuzz/
}
