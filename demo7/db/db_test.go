package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func FuzzAddUser(f *testing.F) {
	f.Add("bernie")

	db, mock, err := sqlmock.New()
	if err != nil {
		f.Fatal(err)
	}
	defer db.Close()

	dbs := NewDatabaseService(db)

	pwd := "supersecretpassword"

	f.Fuzz(func(t *testing.T, username string) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users").WithArgs(username, pwd).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		if err := dbs.AddUser(username, pwd); err != nil {
			t.Errorf("error during processing: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("mock unexpected expectations: %v", err)
		}
	})
}
