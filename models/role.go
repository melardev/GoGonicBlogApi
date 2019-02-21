package models

import "github.com/jinzhu/gorm"

type Role struct {
	gorm.Model
	Name        string
	Description string
	Users       []User `gorm:"many2many:users_roles;"`
}

type UserRole struct {
	User   User `gorm:"foreignkey:UserId"`
	UserId uint
	Role   User `gorm:"foreignkey:RoleId"`
	RoleId uint
}

func (UserRole) TableName() string {
	return "users_roles"
}

func Any(roles []Role, f func(Role) bool) bool {
	for _, role := range roles {
		if f(role) {
			return true
		}
	}
	return false
}
