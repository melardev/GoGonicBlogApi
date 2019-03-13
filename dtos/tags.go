package dtos

import (
	"github.com/melardev/GoGonicBlogApi/models"
)

type CreateTag struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description" binding:"required"`
}

func CreateTagListDto(tags []models.Tag) map[string]interface{} {
	result := map[string]interface{}{}
	var t = make([]interface{}, len(tags))
	for i := 0; i < len(tags); i++ {
		t[i] = CreateTagDto(tags[i])
	}
	result["tags"] = t
	return CreateSuccessDto(result)
}

func CreateTagDto(tag models.Tag) map[string]interface{} {
	return map[string]interface{}{
		"id":          tag.ID,
		"name":        tag.Name,
		"slug":        tag.Slug,
		"description": tag.Description,
	}
}

func CreateTagCreatedDto(tag models.Tag) map[string]interface{} {
	return CreateSuccessWithDtoAndMessageDto(CreateTagDto(tag), "Tag created successfully")
}
