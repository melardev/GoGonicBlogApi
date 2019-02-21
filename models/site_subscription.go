package models

import (
	"github.com/jinzhu/gorm"
)

type SiteSubscription struct {
	gorm.Model
	User   User `gorm:"foreignKey:UserId"`
	UserId uint
}
