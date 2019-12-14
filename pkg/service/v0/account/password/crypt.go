package password

import (
	"golang.org/x/crypto/bcrypt"
)

// CryptHasher uses bcrypt internally to create and compare password hashes
type CryptHasher struct {
}

func NewCryptHasher() *CryptHasher {
	return &CryptHasher{}
}

func (h *CryptHasher) HashPassword(p string) string {
	// Uh oh, assumes no errors for now
	bytes, _ := bcrypt.GenerateFromPassword([]byte(p), 10)
	return string(bytes)
}

func (h *CryptHasher) Compare(hash string, p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(p))
	return err == nil
}
