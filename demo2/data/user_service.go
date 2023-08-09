package data

import (
	"github.com/pkg/errors"

	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

const minEntropyBits = 60

// UserService holds 
type UserService struct {
	users map[string][]byte
}

// NewUserService creates a ready to use user service.
func NewUserService() *UserService {
	return &UserService{
		users: map[string][]byte{},
	}
}

// Add validates a user password and creates a new user. 
func (us *UserService) Add(name, password string) error {
	_, ok := us.users[name]
	if ok {
		return errors.New("user exists already, please log in instead")
	}

	err := passwordvalidator.Validate(password, minEntropyBits)
	if err != nil {
		return errors.Wrap(err, "validate new user password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	us.users[name] = hashedPassword
	return nil
}

// ValidatePassword checks the provided password of an existing user.
func (us *UserService) ValidatePassword(name, providedPwd string) error {
	hashedPassword, ok := us.users[name]
	if !ok {
		return errors.New("user does not exist")
	}
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(providedPwd))
}

// ValidateUser checks the provided username belongs to an existing user. 
func (us *UserService) ValidateUser(name string) error {
	_, ok := us.users[name]
	if !ok {
		return errors.New("user not found")
	}
	return nil
}
