package models

import (
	"github.com/jinzhu/gorm"
)

type Like struct {
	gorm.Model
	Article   Article `gorm:"foreignKey:ArticleId"`
	ArticleId uint
	User      User `gorm:"foreignKey:UserId"`
	UserId    uint
}
