package password

// Hasher defines how a password hasher should behave
type Hasher interface {
	HashPassword(string) string
}
