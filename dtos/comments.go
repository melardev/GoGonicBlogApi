package dtos

import (
	"github.com/gin-gonic/gin"
	"github.com/melardev/api_blog_app/models"
	"net/http"
	"time"
)

type CreateComment struct {
	Content string `form:"content" json:"content" xml:"content"  binding:"required"`
}

func GetCommentDto(comment *models.Comment, includeUser, includeArticle bool) map[string]interface{} {
	result := map[string]interface{}{
		"id":         comment.ID,
		"content":    comment.Content,
		"created_at": comment.CreatedAt.UTC().Format(time.RFC1123),
		"updated_at": comment.UpdatedAt.UTC().Format(time.RFC1123),
	}
	if includeUser == true {
		result["user"] = map[string]interface{}{
			"id":       comment.User.ID,
			"username": comment.User.Username,
		}
	}
	if includeArticle == true {
		result["article"] = map[string]interface{}{
			"id":    comment.Article.ID,
			"slug":  comment.Article.Slug,
			"title": comment.Article.Title,
		}
	}
	return result
}
func CreatedCommentPagedResponse(request *http.Request, comments []models.Comment, page, page_size, count int, bools ...bool) map[string]interface{} {
	var resources = make([]interface{}, len(comments))
	for index, comment := range comments {
		includeUser := false
		if len(bools) > 0 {
			includeUser = bools[0]
		}
		includeArticle := false
		if len(bools) > 1 {
			includeArticle = bools[1]
		}

		resources[index] = GetCommentDto(&comment, includeUser, includeArticle)
	}
	return CreatePagedResponse(request, resources, "comments", page, page_size, count)
}

func CreateCommentDto(comment *models.Comment, includes ...bool) map[string]interface{} {
	includeUser := false
	if len(includes) > 0 {
		includeUser = includes[0]
	}
	includeArticle := false
	if len(includes) > 1 {
		includeArticle = includes[1]
	}

	return GetCommentDto(comment, includeUser, includeArticle)
}

func CreateCommentCreatedDto(comment *models.Comment, includes ...bool) gin.H {
	return CreateSuccessWithDtoAndMessageDto(CreateCommentDto(comment, includes...), "Comment Created successfully")
}
