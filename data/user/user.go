package user

import (
	"fmt"

	"github.com/ian-kent/go-webapp/data"
	"golang.org/x/crypto/bcrypt"
)

// User is a user
type User struct {
	Email    string
	Password []byte
}

// ValidatePassword validates a users password
func (u User) ValidatePassword(password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(password)); err != nil {
		return false, fmt.Errorf("error validating password: %s", err)
	}
	return true, nil
}

// Get gets a user
func Get(email string) (*User, error) {
	if u, ok := data.Storage.Get("user", email); ok {
		return u.(*User), nil
	}
	return nil, nil
}

// Create creates a new user
func Create(email, password string) (*User, error) {
	if _, ok := data.Storage.Get("user", email); ok {
		return nil, fmt.Errorf("email already registered: %s", email)
	}

	b, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	if err != nil {
		return nil, fmt.Errorf("error bcrypting password: %s", err)
	}

	u := &User{
		Email:    email,
		Password: b,
	}

	e := data.Storage.Store("user", email, u)

	return u, e
}
