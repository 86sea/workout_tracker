package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100)"`
	Email     string `gorm:"uniqueIndex;type:varchar(100)"`
	Password  string
	Sets []Set
}

type Set struct {
	ID           uint `gorm:"primaryKey"`
	Date         time.Time
	Reps         uint
	Weight       uint
	ExerciseName string
    ExerciseType string
	UserID       uint `gorm:"foreignKey:ID"`
}
