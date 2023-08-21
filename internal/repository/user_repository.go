package repository

import (
	"workout_tracker/internal/database"
    "workout_tracker/internal/models"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

// CreateUser create a user entry in the user's table
func CreateUser(user *models.User) *gorm.DB {
	return database.DB.Instance.Create(user)
}

// FindUser searches the user's table with the condition given
func FindUser(dest interface{}, conds ...interface{}) *gorm.DB {
	return database.DB.Instance.Model(&models.User{}).Take(dest, conds...)
}

// FindUserByEmail searches the user's table with the email given
func FindUserByEmail(dest interface{}, email string) *gorm.DB {
	return FindUser(dest, "email = ?", email)
}
