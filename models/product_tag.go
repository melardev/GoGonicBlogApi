package models

type ArticleTag struct {
	Tag       User `gorm:"foreignkey:TagId"`
	TagId     uint
	Article   Article `gorm:"foreignkey:ArticleId"`
	ArticleId uint
}

func (ArticleTag) TableName() string {
	return "articles_tags"
}
