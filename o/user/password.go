package user

import (
	"golang.org/x/crypto/bcrypt"
)

type Password string

const LENGTH = 10

func (p Password) GererateHashedPassword() (Password, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(p), LENGTH)
	return Password(hashed), err
}

func (p Password) ComparePassword(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(p), []byte(pwd))
}
