package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"uniqueIndex;type:varchar(255);not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Verified bool   `gorm:"default:false"`
}