package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/melardev/api_blog_app/dtos"
	"github.com/melardev/api_blog_app/infrastructure"
	"github.com/melardev/api_blog_app/models"
	"net/http"
	"os"
	"strings"
)

func EnforceAuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("currentUser")
		if exists && user.(models.User).ID != 0 {
			return
		} else {
			err, exists := c.Get("authErr")
			if exists {
				c.AbortWithStatusJSON(http.StatusForbidden, dtos.CreateDetailedErrorDto("auth_error", err.(error)))
			} else {
				c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("You must be authenticated"))
				c.Abort()
			}
		}
	}
}

func UserLoaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		if bearer != "" {
			jwtParts := strings.Split(bearer, " ")
			if len(jwtParts) == 2 {
				jwtEncoded := jwtParts[1]

				token, err := jwt.Parse(jwtEncoded, func(token *jwt.Token) (interface{}, error) {
					// Theorically we have also to validate the algorithm
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
					}
					secret := []byte(os.Getenv("JWT_SECRET"))
					return secret, nil
				})

				if err != nil {
					println(err.Error())
					return
				}
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					//
					if userId, ok := claims["user_id"]; ok {
						userId = uint(userId.(float64))
						fmt.Printf("[+] Authenticated request, authenticated user id is %d\n", userId)

						var user models.User
						if userId != 0 {
							database := infrastructure.GetDB()
							// We always need the Roles to be loaded to make authorization decisions based on Roles
							database.Preload("Roles").First(&user, userId)
						}

						c.Set("currentUser", user)
						c.Set("currentUserId", user.ID)
					}

				} else {

				}

			}
		}
	}
}

func ShouldBeAuthorOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("currentUser")
		userModel := user.(models.User)
		if exists && (userModel.IsAdmin() || userModel.IsAuthor()) {
			return
		} else {
			// Approach 1
			// c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Permission denied you must be an author or admin"))
			// Prevent any other handler from being called
			// c.Abort()

			// Approach 2
			c.AbortWithStatusJSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Permission denied you must be an author or admin"))
			return
		}
	}
}
