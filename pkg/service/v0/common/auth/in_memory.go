package auth

import (
	"encoding/base64"
	"errors"
	"math/rand"
)

// InMemorySession is an in-memory hashmap holding opaque tokens
type InMemorySession struct {
	db map[string]uint
}

func NewInMemorySession() *InMemorySession {
	var s InMemorySession
	s.db = make(map[string]uint)
	return &s
}

func (s *InMemorySession) GenerateToken(id uint) string {
	b := make([]byte, 10)
	rand.Read(b)
	token := base64.URLEncoding.EncodeToString(b)
	s.db[token] = id
	return token
}

func (s *InMemorySession) FetchID(token string) (uint, error) {
	id, ok := s.db[token]
	if !ok {
		return 0, errors.New("Invalid token")
	}
	return id, nil
}

func (s *InMemorySession) RefreshToken(token string) (ok bool) {
	return true
}
