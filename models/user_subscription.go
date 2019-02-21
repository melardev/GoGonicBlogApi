package models

import "github.com/jinzhu/gorm"

type UserSubscription struct {
	gorm.Model
	Following   User `gorm:"foreignKey:FollowingId"`
	FollowingId uint
	Follower    User `gorm:"foreignKey:FollowerId"`
	FollowerId  uint
}

func (UserSubscription) TableName() string {
	return "users_subscriptions"
}
