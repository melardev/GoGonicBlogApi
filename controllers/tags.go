package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/melardev/api_blog_app/dtos"
	"github.com/melardev/api_blog_app/services"
	"net/http"
)

func RegisterTagRoutes(router *gin.RouterGroup) {
	router.GET("", TagList)
}
func TagList(c *gin.Context) {
	tags, err := services.FetchAllTags()
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	c.JSON(http.StatusOK, dtos.CreateTagListDto(tags))
}
