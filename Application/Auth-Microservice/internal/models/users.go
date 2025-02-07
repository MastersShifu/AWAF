package models

type Users struct {
	ID       string `gorm:"PRIMARY_KEY"`
	Name     string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"not null"`
	Version  int    `gorm:"version"`
}

func (Users) TableName() string {
	return "users"
}
