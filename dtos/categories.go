package dtos

import "github.com/melardev/api_blog_app/models"

func CreateCategoryListDto(categories []models.Category) map[string]interface{} {
	result := map[string]interface{}{}
	var t = make([]interface{}, len(categories))
	for i := 0; i < len(categories); i++ {
		t[i] = CreateCategoryDto(categories[i])
	}
	result["categories"] = t
	return CreateSuccessDto(result)
}

func CreateCategoryDto(tag models.Category) map[string]interface{} {
	return map[string]interface{}{
		"id":          tag.ID,
		"name":        tag.Name,
		"slug":        tag.Slug,
		"description": tag.Description,
	}
}

func CreateCategoryCreatedDto(category models.Category) map[string]interface{} {
	return CreateSuccessWithDtoAndMessageDto(CreateCategoryDto(category), "Category created successfully")
}
