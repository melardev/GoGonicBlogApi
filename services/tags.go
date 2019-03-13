package services

import (
	"github.com/melardev/GoGonicBlogApi/infrastructure"
	"github.com/melardev/GoGonicBlogApi/models"
)

func FetchAllTags() ([]models.Tag, error) {
	database := infrastructure.GetDB()
	var tags []models.Tag
	err := database.Find(&tags).Error
	return tags, err
}
