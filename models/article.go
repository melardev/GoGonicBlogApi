package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Slug          string `gorm:"unique_index"`
	Title         string
	Description   string `gorm:"size:2048"`
	Body          string `gorm:"size:2048"`
	User          User
	UserId        uint
	Tags          []Tag      `gorm:"many2many:articles_tags;"`
	Categories    []Category `gorm:"many2many:articles_categories;"`
	Comments      []Comment  `gorm:"foreignKey:ArticleId"`
	CommentsCount int        `gorm:"-"`
}

func (a *Article) BeforeSave() (err error) {
	a.Slug = slug.Make(a.Title)
	return
}
