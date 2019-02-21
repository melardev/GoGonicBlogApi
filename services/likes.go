package services

import (
	"github.com/melardev/api_blog_app/infrastructure"
	"github.com/melardev/api_blog_app/models"
)

func LikesCount(article *models.Article) uint {
	database := infrastructure.GetDB()
	var count uint
	database.Model(&models.Like{}).Where(models.Like{
		ArticleId: article.ID,
	}).Count(&count)
	return count
}
