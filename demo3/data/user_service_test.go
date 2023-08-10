package data_test

import (
	"testing"

	"github.com/addetz/secure-code-go/demo3/data"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	t.Run("insufficient password", func(t *testing.T) {
		us := data.NewUserService()
		err := us.Add("new", "test")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "validate new user password: insecure password")
	})
	t.Run("successful add", func(t *testing.T) {
		name := "user1"
		password := "test-horse-pen-clam"
		us := data.NewUserService()
		err := us.Add(name, password)
		assert.Nil(t, err)
	})
	t.Run("duplicate user", func(t *testing.T) {
		name := "user1"
		password := "test-horse-pen-clam"
		us := data.NewUserService()
		err := us.Add(name, password)
		assert.Nil(t, err)
		err = us.Add(name, password)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "user exists already")
	})
}

func TestValidate(t *testing.T) {
	t.Run("successful validate", func(t *testing.T) {
		name := "user1"
		password := "test-horse-pen-clam"
		us := data.NewUserService()
		err := us.Add(name, password)
		assert.Nil(t, err)
		err = us.ValidatePassword(name, password)
		assert.Nil(t, err)
	})
	t.Run("failed validate", func(t *testing.T) {
		name := "user1"
		password := "test-horse-pen-clam"
		us := data.NewUserService()
		err := us.Add(name, password)
		assert.Nil(t, err)
		err = us.ValidatePassword(name, "garbage-password")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "hashedPassword is not the hash of the given password")
	})
	t.Run("inexistent user", func(t *testing.T) {
		name := "user1"
		us := data.NewUserService()
		err := us.ValidatePassword(name, "garbage-password")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "user does not exist")
	})
}
