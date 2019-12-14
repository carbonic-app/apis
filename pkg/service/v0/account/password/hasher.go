package password

// Hasher defines how a password hasher should behave
type Hasher interface {
	HashPassword(string) string
	Compare(hash string, p string) bool
}
