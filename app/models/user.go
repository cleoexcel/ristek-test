package models

type User struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}
