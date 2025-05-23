package auth

import (
	"errors"
	"fmt"
)

var errWrongCredentials = errors.New("user with specified credentials not found")

func errUserAlreadyExists(login string) error {
	msg := fmt.Sprintf("user with login %s already exists", login)
	return errors.New(msg)
}
