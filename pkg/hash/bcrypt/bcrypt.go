package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptProvider struct{}

func NewBcryptProvider() *BcryptProvider {
	return new(BcryptProvider)
}

func (p *BcryptProvider) Make(input string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	return string(b)
}

func (p *BcryptProvider) Compare(hashedInput, input string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedInput), []byte(input))
}
