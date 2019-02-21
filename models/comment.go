package models

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	Content          string  `gorm:"size:2048"`
	Article          Article `gorm:"foreignKey:ArticleId"`
	ArticleId        uint
	User             User `gorm:"foreignKey:UserId"`
	UserId           uint
	RepliedComment   *Comment `gorm:"foreignkey:RepliedCommentID"`
	RepliedCommentID *uint
	Replies          []*Comment `gorm:"foreignkey:ID"`
}

func (comment *Comment) BeforeSave() (err error) {
	if comment.RepliedComment != nil {
		comment.RepliedCommentID = &comment.RepliedComment.ID
	}
	return
}
