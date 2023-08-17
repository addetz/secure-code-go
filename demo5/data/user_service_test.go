package data_test

import (
	"errors"
	"testing"

	"github.com/addetz/secure-code-go/demo4/data"
	"github.com/addetz/secure-code-go/demo4/db"
	"github.com/addetz/secure-code-go/demo4/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAdd(t *testing.T) {
	t.Run("insufficient password", func(t *testing.T) {
		name := "user1"
		mockDB := new(mocks.DatabaseServiceMock)
		us := data.NewUserService(mockDB)
		mockDB.On("GetUser", name).Return(nil, errors.New("no user exists"))
		err := us.Add(name, "test")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "validate new user password: insecure password")
	})
	t.Run("successful add", func(t *testing.T) {
		name := "user1"
		password := "test-horse-pen-clam"
		mockDB := new(mocks.DatabaseServiceMock)
		us := data.NewUserService(mockDB)
		mockDB.On("GetUser", name).Return(nil, errors.New("no user exists"))
		mockDB.On("AddUser", name, mock.AnythingOfType("string")).Return(nil)
		err := us.Add(name, password)
		assert.Nil(t, err)
	})
	t.Run("duplicate user", func(t *testing.T) {
		name := "user1"
		password := "test-horse-pen-clam"
		mockDB := new(mocks.DatabaseServiceMock)
		us := data.NewUserService(mockDB)
		mockDB.On("GetUser", name).Return(nil, nil)
		err := us.Add(name, password)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "user exists already")
	})
}

func TestValidate(t *testing.T) {
	t.Run("successful validate", func(t *testing.T) {
		name := "user1"
		password := "test-horse-pen-clam"
		expected, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		mockDB := new(mocks.DatabaseServiceMock)
		us := data.NewUserService(mockDB)
		mockDB.On("GetUser", name).Return(&db.User{
			Username: name,
			Pwd:      string(expected),
		}, nil)
		err := us.ValidatePassword(name, password)
		assert.Nil(t, err)
	})
	t.Run("failed validate", func(t *testing.T) {
		name := "user1"
		password := "test-horse-pen-clam"
		expected, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		mockDB := new(mocks.DatabaseServiceMock)
		mockDB.On("GetUser", name).Return(&db.User{
			Username: name,
			Pwd:      string(expected),
		}, nil)
		us := data.NewUserService(mockDB)
		err := us.ValidatePassword(name, "garbage-password")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "hashedPassword is not the hash of the given password")
	})
	t.Run("inexistent user", func(t *testing.T) {
		name := "user1"
		mockDB := new(mocks.DatabaseServiceMock)
		us := data.NewUserService(mockDB)
		mockDB.On("GetUser", name).Return(nil, errors.New("no user exists"))
		err := us.ValidatePassword(name, "garbage-password")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "user does not exist")
	})
}
