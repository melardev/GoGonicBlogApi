package services

import (
	"github.com/melardev/GoGonicBlogApi/infrastructure"
	"github.com/melardev/GoGonicBlogApi/models"
)

func LikesCount(article *models.Article) uint {
	database := infrastructure.GetDB()
	var count uint
	database.Model(&models.Like{}).Where(models.Like{
		ArticleId: article.ID,
	}).Count(&count)
	return count
}
