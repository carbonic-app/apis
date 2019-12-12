package auth

// Session helps other services fetch IDs out of tokens
type Session interface {
	GenerateToken(id uint) string
	FetchID(token string) (uint, error)
	RefreshToken(token string) (ok bool)
}
