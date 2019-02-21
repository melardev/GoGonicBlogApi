package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type User struct {
	gorm.Model
	Username  string  `gorm:"column:username"`
	FirstName string  `gorm:"varchar(255);not null"`
	LastName  string  `gorm:"varchar(255);not null"`
	Email     string  `gorm:"column:email;unique_index"`
	Bio       string  `gorm:"column:bio;size:1024"`
	Image     *string `gorm:"column:image"`
	Password  string  `gorm:"column:password;not null"`

	Roles []Role `gorm:"many2many:users_roles;"`

	Following []*User `gorm:"many2many:users_subscriptions;association_jointable_foreignkey:follower_id"`
	Followers []*User `gorm:"many2many:users_subscriptions;association_jointable_foreignkey:following_id"`
}

func (user *User) BeforeSave(db *gorm.DB) (err error) {
	if len(user.Roles) == 0 {
		role := Role{}
		userRole := Role{}
		db.Model(&role).Where("name = ?", "ROLE_USER").First(&userRole)
		user.Roles = append(user.Roles, userRole)
	}
	return
}

func (user *User) IsValidPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(user.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

// Generate JWT token associated to this user
func (user *User) GenerateJwtToken() string {
	// jwt.New(jwt.GetSigningMethod("HS512"))
	jwt_token := jwt.New(jwt.SigningMethodHS512)

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	jwt_token.Claims = jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"roles":    roles,
		"exp":      time.Now().Add(time.Hour * 24 * 90).Unix(),
	}
	// Sign and get the complete encoded token as a string
	token, _ := jwt_token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return token
}
func (user *User) IsAuthor() bool {
	for _, role := range user.Roles {
		if role.Name == "ROLE_AUTHOR" {
			return true
		}
	}
	return false
}

func (user *User) IsAdmin() bool {
	for _, role := range user.Roles {
		if role.Name == "ROLE_ADMIN" {
			return true
		}
	}
	return false
}

func (user *User) IsNotAdmin() bool {
	return !user.IsAdmin()
}
