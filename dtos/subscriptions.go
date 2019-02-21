package dtos

import (
	"github.com/melardev/api_blog_app/models"
	"net/http"
)

func CreateUserSubscriptionPageResponse(request *http.Request, users []models.User, resourceName string, page, pageSize, count int) map[string]interface{} {
	var resources = make([]interface{}, len(users))
	for index, user := range users {
		resources[index] = GetUserBasicInfo(user)
	}
	return CreatePagedResponse(request, resources, resourceName, page, pageSize, count)
}

func CreateFollowersFollowingPageResponse(request *http.Request, following []models.User, followers []models.User, page, pageSize, count int) map[string]interface{} {
	var followersResponse = make([]interface{}, len(followers))
	for index, user := range followers {
		followersResponse[index] = GetUserBasicInfo(user)
	}

	var followingResponse = make([]interface{}, len(following))
	for index, user := range following {
		followingResponse[index] = GetUserBasicInfo(user)
	}
	response := CreatePageMeta(request, len(followers)+len(following), page, pageSize, count)
	response["followers"] = followersResponse
	response["following"] = followingResponse
	return response
}
