package database

import (
	"main/models"
)

func AutoMigrate() {
	db.AutoMigrate(&models.User{}, &models.Article{}, &models.Comment{})
}
