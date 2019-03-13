package services

import (
	"github.com/melardev/GoGonicBlogApi/infrastructure"
	"github.com/melardev/GoGonicBlogApi/models"
)

func FetchAllCategories() ([]models.Category, error) {
	database := infrastructure.GetDB()
	var categories []models.Category
	err := database.Find(&categories).Error
	return categories, err
}
