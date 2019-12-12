package common

// Auth helps other services fetch IDs out of tokens
type Auth interface {
	GenerateToken(id uint) string
	FetchID(token string) uint
	RefreshToken(token string) bool
}
