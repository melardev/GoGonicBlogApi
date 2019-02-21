package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Category struct {
	gorm.Model
	Name        string `gorm:"unique_index"`
	Slug        string `gorm:"unique_index"`
	Description string
	Articles    []Article `gorm:"many2many:articles_categories;"`
	IsNewRecord bool      `gorm:"-;default:false"` // Virtual Field, so it is not persisted in the Db. This is used in FirstOrCreate()
}

func (a *Category) BeforeSave() (err error) {
	a.Slug = slug.Make(a.Name)
	return
}
