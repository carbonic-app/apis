package password

// PlaintextHasher is the simplest hasher in that it doesn't hash the password
type PlaintextHasher struct {
}

func NewPlaintextHasher() *PlaintextHasher {
	return &PlaintextHasher{}
}

func (h *PlaintextHasher) HashPassword(p string) string {
	return p
}

func (h *PlaintextHasher) Compare(hash string, p string) bool {
	return hash == p
}
