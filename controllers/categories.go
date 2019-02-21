package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/melardev/api_blog_app/dtos"
	"github.com/melardev/api_blog_app/services"
	"net/http"
)

func RegisterCategoryRoutes(router *gin.RouterGroup) {
	router.GET("", ListCategories)
}
func ListCategories(c *gin.Context) {
	categories, err := services.FetchAllCategories()
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	c.JSON(http.StatusOK, dtos.CreateCategoryListDto(categories))
}
