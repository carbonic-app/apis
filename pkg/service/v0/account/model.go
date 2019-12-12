package account

type User struct {
	ID           uint `gorm:"primary_key"`
	Username     string
	PasswordHash string
}
