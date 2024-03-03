package util

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct {
}

func (b *Bcrypt) HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword, err
}

func (b *Bcrypt) ComparePassword(password string, compare []byte) error {
	err := bcrypt.CompareHashAndPassword(compare, []byte(password))
	return err
}
