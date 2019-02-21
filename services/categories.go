package services

import (
	"github.com/melardev/api_blog_app/infrastructure"
	"github.com/melardev/api_blog_app/models"
)

func FetchAllCategories() ([]models.Category, error) {
	database := infrastructure.GetDB()
	var categories []models.Category
	err := database.Find(&categories).Error
	return categories, err
}
