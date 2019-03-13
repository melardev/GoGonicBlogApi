package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/melardev/GoGonicBlogApi/dtos"
	"github.com/melardev/GoGonicBlogApi/infrastructure"
	"github.com/melardev/GoGonicBlogApi/middlewares"
	"github.com/melardev/GoGonicBlogApi/models"
	"github.com/melardev/GoGonicBlogApi/services"
	"net/http"
	"strconv"
)

func RegisterLikeRoutes(router *gin.RouterGroup) {
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.POST("/articles/:slug/likes", LikeArticle)
		router.DELETE("/articles/:slug/likes", DislikeArticle)
		router.GET("/likes", MyLikes)
	}
}

func MyLikes(c *gin.Context) {
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")
	currentUserId := c.MustGet("currentUserId").(uint)
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	var result []models.Like
	var ids []uint
	database := infrastructure.GetDB()

	// Approach 1
	/*
		rows, err := database.Table("likes").Select("article_id").Where("user_id = ?", currentUserId).Rows()
		for rows.Next() {
			var articleId interface{}
			rows.Scan(&articleId)
			// then append articleId on ids []uint
			fmt.Println(res)
		}
	*/
	// Approach 2, getting array of uint(ids) populated
	// err = database.Table("likes").Select("article_id").Where("user_id = ?", currentUserId).Scan(&ids).Error

	// Approach 3, getting array of Like with ArticleId field populated
	// database.Select([]string{"article_id", "user_id"}).Where("likes.user_id = ?", currentUserId).Find(&result)
	// Approach 4, getting array of Like with ArticleId field populated
	// database.Select("article_id").Where(&models.Like{UserId: currentUserId}).Find(&result)

	// Approach 5
	database.Select("article_id").Where("likes.user_id = ?", currentUserId).Find(&result).Pluck("article_id", &ids)

	var articles []models.Article
	var likedArticles = 0
	database.Table("articles").Where("id in (?)", ids).Count(&likedArticles)
	database.Where("id in (?)", ids).
		Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).
		Preload("User").Preload("Tags").Preload("Categories").Find(&articles)
	// For some reason if we use Select() it crashes when executing the loop below, TODO: fix it
	// .Select([]string{"title", "slug", "user_id", "description"})
	// I may also get ride of comments_count tags and categories, this endpoint should only return id, title, slug, user.id and user.username

	// TODO: Performance to be improved, at this state we are making a single SQL query PER article
	// we should make a single SQL Query for all the articles
	for index := range articles {
		articles[index].CommentsCount = database.Model(&articles[index]).Association("Comments").Count()
	}

	c.JSON(http.StatusOK, dtos.CreatedArticlePagedResponse(c.Request, articles, page, pageSize, likedArticles))
}

func LikeArticle(c *gin.Context) {
	slug := c.Param("slug")

	database := infrastructure.GetDB()
	var article models.Article
	err := database.Model(&models.Article{}).Where("slug = ?", slug).Select([]string{"id", "title"}).First(&article).Error
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("articles", err))
		return
	}

	user := c.MustGet("currentUser").(models.User)

	if !services.IsLikedBy(&article, user) {
		like := models.Like{
			ArticleId: article.ID,
			UserId:    user.ID,
		}
		// article.LikeArticle(user)
		if err := services.CreateOne(&like); err != nil {
			c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database", err))
			return
		}

		c.JSON(http.StatusOK, dtos.GetSuccessDto(fmt.Sprintf("You liked the article \"%v\" successfully", article.Title)))
	} else {
		c.JSON(http.StatusForbidden, dtos.GetErrorDto("You have already liked this article"))
	}

}

func DislikeArticle(c *gin.Context) {
	slug := c.Param("slug")
	database := infrastructure.GetDB()
	var result struct {
		Id string
	}
	err := database.Table("articles").Select("id").Where("slug = ?", slug).Scan(&result).Error
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("articles", err))
		return
	}
	user := c.MustGet("currentUser").(models.User)
	var like models.Like
	err = database.Model(models.Like{}).Where("user_id = ? AND article_id = ?", user.ID, result.Id).First(&like).Error
	if err == nil {
		// err = article.Unlike(user)
		database.Delete(&like)
		c.JSON(http.StatusOK, dtos.GetSuccessDto("Article disliked successfully"))
	} else {
		c.JSON(http.StatusForbidden, dtos.GetSuccessDto("You were not liking this article, so you can not perform this operation"))
	}
}
