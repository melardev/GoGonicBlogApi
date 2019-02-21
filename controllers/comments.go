package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/melardev/api_blog_app/dtos"
	"github.com/melardev/api_blog_app/infrastructure"
	"github.com/melardev/api_blog_app/middlewares"
	"github.com/melardev/api_blog_app/models"
	"github.com/melardev/api_blog_app/services"
	"net/http"
	"strconv"
)

func RegisterCommentRoutes(router *gin.RouterGroup) {
	router.GET("/articles/:slug/comments", ListComments)
	router.GET("/articles/:slug/comments/:id", ShowComment)
	router.GET("/comments/:id", ShowComment)

	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.POST("/articles/:slug/comments", CreateComment)
		router.DELETE("/comments/:id", DeleteComment)
		router.DELETE("/articles/:slug/comments/:id", DeleteComment)
	}

}

func ListComments(c *gin.Context) {
	slug := c.Param("slug")
	var article models.Article

	database := infrastructure.GetDB()
	articleId := -1
	err := database.Model(&article).Where(&models.Article{Slug: slug}).Select("id").Row().Scan(&articleId)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("comments", err))
		return
	}
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	var comments []models.Comment
	totalCommentCount := 0
	database.Model(&comments).Where(&models.Comment{ArticleId: uint(articleId)}).Count(&totalCommentCount)
	database.Offset((page - 1) * pageSize).Limit(pageSize).Where(&models.Comment{ArticleId: uint(articleId)}).Preload("User").Find(&comments)

	c.JSON(http.StatusOK, dtos.CreatedCommentPagedResponse(c.Request, comments, page, pageSize, totalCommentCount, true, false))
}

func CreateComment(c *gin.Context) {
	slug := c.Param("slug")
	articleId, err := services.FetchArticleId(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("comment", err))
		return
	}

	var json dtos.CreateComment
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	comment := models.Comment{
		Content:   json.Content,
		ArticleId: articleId,
		User:      c.MustGet("currentUser").(models.User),
	}

	if err := services.SaveOne(&comment); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("Database error", err))
		return
	}

	c.JSON(http.StatusOK, dtos.CreateCommentCreatedDto(&comment))
}

func ShowComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dtos.GetErrorDto("You must provide a valid comment id"))
	}
	comment := services.FetchCommentById(id, true, true)
	c.JSON(http.StatusOK, dtos.CreateCommentDto(&comment, true, true))
}

func DeleteComment(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)
	database := infrastructure.GetDB()
	var comment models.Comment
	err = database.Select([]string{"id", "user_id"}).Find(&comment, id).Error
	if err != nil || comment.ID == 0 {
		// the comment.ID == is redundant, but shows the other way of checking but it is less readable
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("db_error", err))
	} else if currentUser.ID == comment.UserId || currentUser.IsAdmin() {
		err = database.Delete(&comment).Error
		if err != nil {
			c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("comment", err))
			return
		}
		c.JSON(http.StatusOK, dtos.GetSuccessDto("Comment Deleted successfully"))
	} else {
		c.JSON(http.StatusForbidden, dtos.GetErrorDto("You have to be admin or the owner of this comment to delete it"))
	}
}
