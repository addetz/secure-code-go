package data

import (
	"github.com/addetz/secure-code-go/demo4/db"
	"github.com/pkg/errors"

	passwordvalidator "github.com/wagslane/go-password-validator"
	"golang.org/x/crypto/bcrypt"
)

const minEntropyBits = 60

type User struct {
	Username, Password string
}

// UserService holds
type UserService struct {
	dbService db.DatabaseService
}

// NewUserService creates a ready to use user service.
func NewUserService(dbService db.DatabaseService) *UserService {
	return &UserService{
		dbService: dbService,
	}
}

// Add validates a user password and creates a new user.
func (us *UserService) Add(name, password string) error {
	if _, err := us.dbService.GetUser(name); err == nil {
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
	return us.dbService.AddUser(name, string(hashedPassword))
}

// ValidatePassword checks the provided password of an existing user.
func (us *UserService) ValidatePassword(name, providedPwd string) error {
	user, err := us.dbService.GetUser(name)
	if err != nil {
		return errors.New("user does not exist")
	}
	return bcrypt.CompareHashAndPassword([]byte(user.Pwd), []byte(providedPwd))
}

// ValidateUser checks the provided username belongs to an existing user.
func (us *UserService) ValidateUser(name string) error {
	if _, err := us.dbService.GetUser(name); err != nil {
		return errors.New("user not found")
	}
	return nil
}
