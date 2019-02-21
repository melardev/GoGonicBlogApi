package services

import (
	"github.com/melardev/api_blog_app/infrastructure"
	"github.com/melardev/api_blog_app/models"
)

func FetchAllTags() ([]models.Tag, error) {
	database := infrastructure.GetDB()
	var tags []models.Tag
	err := database.Find(&tags).Error
	return tags, err
}
