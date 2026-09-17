package models

type Users struct {
	Base
	Name     string `gorm:"size:100;not null"`
	Email    string `gorm:"size:255;uniqueIndex;not null"`
	Password string `gorm:"size:255;not null"`
}

func (Users) TableName() string {
	return "users"
}
