package controllers

// import "C"
import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/melardev/api_blog_app/dtos"
	"github.com/melardev/api_blog_app/infrastructure"
	"github.com/melardev/api_blog_app/middlewares"
	"github.com/melardev/api_blog_app/models"
	"github.com/melardev/api_blog_app/services"
	"net/http"
	"strconv"
)

func RegisterArticleRoutes(router *gin.RouterGroup) {
	router.GET("/", ListArticles)
	router.GET("/:slug", ShowArticle)

	router.Use(middlewares.EnforceAuthenticatedMiddleware(), middlewares.ShouldBeAuthorOrAdmin())
	{
		router.POST("/", CreateArticle)
		// router.PUT("/:slug", UpdateArticle)
		router.DELETE("/:slug", DeleteArticle)
	}
}

func ListArticles(c *gin.Context) {

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

	articles, modelCount, err := services.FetchArticlesPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("articles", errors.New("invalid param")))
		return
	}

	c.JSON(http.StatusOK, dtos.CreatedArticlePagedResponse(c.Request, articles, page, pageSize, modelCount))
}

func ShowArticle(c *gin.Context) {
	slugParam := c.Param("slug")
	if slugParam == "feed" {
		ArticleFeed(c)
		return
	}
	article, err := services.FetchArticleDetails(&models.Article{Slug: slugParam}, true)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("db_errors", err))
		return
	}

	c.JSON(http.StatusOK, dtos.GetArticleDetailsDto(&article))
}

func CreateArticle(c *gin.Context) {

	var json dtos.CreateArticle
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	database := infrastructure.GetDB()
	tags := make([]models.Tag, len(json.Tags))
	categories := make([]models.Category, len(json.Categories))

	for index, tag := range json.Tags {
		database.Where(&models.Tag{Slug: slug.Make(tag.Name)}).
			Attrs(models.Tag{Name: tag.Name, Description: tag.Description}).
			FirstOrCreate(&tags[index])
	}

	for index, _ := range json.Categories {
		database.Where(&models.Category{Slug: slug.Make(json.Categories[index].Name)}).
			Attrs(models.Category{Name: json.Categories[index].Name, Description: json.Categories[index].Description}).
			FirstOrCreate(&categories[index])
	}

	article := models.Article{
		Title:       json.Title,
		Description: json.Description,
		Body:        json.Body,
		Tags:        tags,
		Categories:  categories,
		UserId:      c.MustGet("currentUserId").(uint),
		User:        c.MustGet("currentUser").(models.User),
	}

	if err := services.CreateOne(&article); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	c.JSON(http.StatusOK, dtos.GetArticleDetailsDto(&article))

}

func DeleteArticle(c *gin.Context) {
	slugParam := c.Param("slugParam")
	user := c.MustGet("currentUser").(models.User)
	err := services.DeleteArticleIfOwnerOrAdmin(&user, &models.Article{Slug: slugParam})
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}
	c.JSON(http.StatusOK, dtos.CreateSuccessWithMessageDto("Article Deleted successfully"))
}

func ArticleFeed(c *gin.Context) {
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

	user := c.MustGet("currentUser").(models.User)

	articles, totalArticlesCount, err := GetFeed(&user, page, pageSize)
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	c.JSON(http.StatusOK, dtos.CreatedArticlePagedResponse(c.Request, articles, page, pageSize, totalArticlesCount))

}

func GetFeed(user *models.User, limit, offset int) ([]models.Article, int, error) {
	database := infrastructure.GetDB()
	var articles []models.Article
	var totalArticlesCount int

	tx := database.Begin()

	followingUserIds := services.GetFollowingIds(user)

	tx.Where("user_id in (?)", followingUserIds).Order("created_at desc").
		Offset(offset).Limit(limit).
		Find(&articles)

	comments := make([]int, len(articles))

	for index, _ := range articles {
		tx.Model(&articles[index]).Related(&articles[index].User, "User")
		tx.Model(&articles[index].User).Related(&articles[index].User)
		tx.Model(&articles[index]).Related(&articles[index].Tags, "Tags")
		comments[index] = tx.Model(&articles[index]).Association("Comments").Count()
	}
	err := tx.Commit().Error
	return articles, totalArticlesCount, err
}
