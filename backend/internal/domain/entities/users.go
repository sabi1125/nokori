package entities

type Users struct {
	UserID       string `json:"id" gorm:"column:user_id"`
	FirstName    string `json:"first_name" gorm:"column:first_name"`
	LastName     string `json:"last_name" gorm:"column:last_name"`
	Email        string `json:"email" gorm:"column:email"`
	PasswordHash string `json:"password_hash" gorm:"column:password_hash"`
	Verified     bool   `json:"verified" gorm:"column:verified"`
}
