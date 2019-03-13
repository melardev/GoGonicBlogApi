package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/melardev/GoGonicBlogApi/infrastructure"
	"github.com/melardev/GoGonicBlogApi/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// What's bcrypt? https://en.wikipedia.org/wiki/Bcrypt
// Golang bcrypt doc: https://godoc.org/golang.org/x/crypto/bcrypt
// You can change the value in bcrypt.DefaultCost to adjust the security index.
// 	err := userModel.setPassword("password0")
func SetPassword(user *models.User, password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	user.Password = string(passwordHash)
	return nil
}

func FindOneUser(condition interface{}) (models.User, error) {
	database := infrastructure.GetDB()
	var user models.User
	err := database.Where(condition).Preload("Roles").First(&user).Error
	return user, err
}

func CreateOne(data interface{}) error {
	database := infrastructure.GetDB()
	err := database.Create(data).Error
	return err
}

func SaveOne(data interface{}) error {
	database := infrastructure.GetDB()
	// test 1 to 1
	// tx1 := database.Begin()
	err := database.Save(data).Error
	// tx1.Commit()
	return err
}

func UpdateUser(user *models.User, data interface{}) error {
	database := infrastructure.GetDB()
	err := database.Model(user).Update(data).Error
	return err
}

func following(user *models.User, v models.User) error {
	database := infrastructure.GetDB()
	var follow models.UserSubscription
	err := database.FirstOrCreate(&follow, &models.UserSubscription{
		FollowingId: v.ID,
		FollowerId:  user.ID,
	}).Error
	return err
}

func IsFollowing(follower *models.User, following models.User) bool {
	database := infrastructure.GetDB()
	var follow models.UserSubscription
	database.Where(models.UserSubscription{
		FollowingId: following.ID,
		FollowerId:  follower.ID,
	}).First(&follow)
	return follow.ID != 0
}

func unFollowing(follower *models.User, following models.User) error {
	database := infrastructure.GetDB()
	err := database.Where(models.UserSubscription{
		FollowingId: following.ID,
		FollowerId:  follower.ID,
	}).Delete(models.UserSubscription{}).Error
	return err
}

func GetFollowings(user *models.User) []models.User {
	database := infrastructure.GetDB()
	tx := database.Begin()
	var follows []models.UserSubscription
	var followings []models.User
	tx.Where(models.UserSubscription{
		FollowerId: user.ID,
	}).Find(&follows)
	for _, follow := range follows {
		var userModel models.User
		tx.Model(&follow).Related(&userModel, "Following")
		followings = append(followings, userModel)
	}
	tx.Commit()
	return followings
}

func GetFollowingIds(user *models.User) []uint {
	database := infrastructure.GetDB()
	tx := database.Begin()
	var users []models.UserSubscription
	var followingIds []uint
	tx.Where(models.UserSubscription{
		FollowerId: user.ID,
	}).Find(&users).Pluck("following_id", &followingIds)

	tx.Commit()

	return followingIds
}

// Generate JWT token associated to this user
func GenerateJwtToken(user *models.User) string {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS512"))

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	jwt_token.Claims = jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"roles":    roles,
		"exp":      time.Now().Add(time.Hour * 24 * 90).Unix(),
	}
	// Sign and get the complete encoded token as a string
	token, _ := jwt_token.SignedString([]byte("JWT_SUPER_SECRET"))
	return token
}
