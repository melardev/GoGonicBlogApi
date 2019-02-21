package models

type ArticleCategory struct {
	Category   User `gorm:"foreignkey:CategoryId"`
	CategoryId uint
	Article    Article `gorm:"foreignkey:ArticleId"`
	ArticleId  uint
}

func (ArticleCategory) TableName() string {
	return "articles_categories"
}
