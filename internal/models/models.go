package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(100)"`
	Email    string `gorm:"uniqueIndex;type:varchar(100)"`
	Password string
}

