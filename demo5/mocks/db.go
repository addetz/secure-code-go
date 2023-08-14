package mocks

import (
	"github.com/addetz/secure-code-go/demo4/db"
	"github.com/stretchr/testify/mock"
)

type DatabaseServiceMock struct {
	mock.Mock
}

func (m *DatabaseServiceMock) AddUser(username, pwd string) error {
	args := m.Called(username, pwd)
	return args.Error(0)
}

func (m *DatabaseServiceMock) GetUser(username string) (*db.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	userArg := args.Get(0).(*db.User)
	return userArg, args.Error(1)
}

func (m *DatabaseServiceMock) AddNote(id, username, text string) error {
	args := m.Called(id, username, text)
	return args.Error(0)
}

func (m *DatabaseServiceMock) GetUserNotes(username string) ([]db.Note, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	notesArg := args.Get(0).([]db.Note)
	return notesArg, args.Error(1)
}
