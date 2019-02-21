package dtos

import (
	"github.com/melardev/api_blog_app/models"
	"net/http"
	"time"
)

type CreateArticle struct {
	Title       string `form:"title" json:"title" xml:"title"  binding:"required"`
	Description string `form:"description" json:"description" xml:"description" binding:"required"`
	Body        string `form:"body" json:"body" xml:"body" binding:"required"`
	Tags        []struct {
		Name        string `form:"name" json:"name" binding:"exists,alphanum,min=4,max=255"`
		Description string `form:"description" json:"description" binding:"exists"`
	} `json:"tags"`
	Categories []struct {
		Name        string `form:"name" json:"name" binding:"exists,alphanum,min=4,max=255"`
		Description string `form:"email" json:"mail" binding:"exists"`
	} `json:"categories"`
}

func GetArticleDto(article *models.Article) map[string]interface{} {
	tagsCount := len(article.Tags)
	var tags = make([]map[string]interface{}, tagsCount)
	var categories = make([]map[string]interface{}, len(article.Categories))
	for index, tag := range article.Tags {
		tags[index] = map[string]interface{}{
			"id":   tag.ID,
			"name": tag.Name,
		}
	}

	for index, category := range article.Categories {
		categories[index] = map[string]interface{}{
			"id":   category.ID,
			"name": category.Name,
		}
	}

	result := map[string]interface{}{
		"id":          article.ID,
		"title":       article.Title,
		"slug":        article.Slug,
		"description": article.Description,
		"user": map[string]interface{}{
			"id":       article.User.ID,
			"username": article.User.Username,
		},
		"tags":       tags,
		"categories": categories,
		"created_at": article.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
		"updated_at": article.UpdatedAt.UTC().Format(time.RFC3339Nano),
	}
	if article.CommentsCount >= 0 {
		result["comments_count"] = article.CommentsCount
	}
	return result
}

func GetArticleDetailsDto(article *models.Article) map[string]interface{} {
	result := GetArticleDto(article)
	result["body"] = article.Body
	comments := make([]map[string]interface{}, len(article.Comments))
	for i := 0; i < len(article.Comments); i++ {
		comments[i] = GetCommentDto(&article.Comments[i], true, false)
	}

	result["comments"] = comments
	return result
}

func CreatedArticlePagedResponse(request *http.Request, articles []models.Article, page, page_size, count int) interface{} {
	var resources = make([]interface{}, len(articles))
	for i := 0; i < len(articles); i++ {
		resources[i] = GetArticleDto(&articles[i])
	}
	return CreatePagedResponse(request, resources, "articles", page, page_size, count)
}
