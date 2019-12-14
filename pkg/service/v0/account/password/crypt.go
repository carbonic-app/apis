package password

import (
	"golang.org/x/crypto/bcrypt"
)

// CryptHasher uses bcrypt internally to create and compare password hashes
type CryptHasher struct {
	salt string
}

func NewCryptHasher(salt string) *CryptHasher {
	return &CryptHasher{salt: salt}
}

func (h *CryptHasher) HashPassword(p string) string {
	// Uh oh, assumes no errors for now
	bytes, _ := bcrypt.GenerateFromPassword([]byte(p+h.salt), 10)
	return string(bytes)
}

func (h *CryptHasher) Compare(hash string, p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(p+h.salt))
	return err == nil
}
