package services

import (
	"github.com/melardev/GoGonicBlogApi/infrastructure"
	"github.com/melardev/GoGonicBlogApi/models"
)

func IsLikedBy(article *models.Article, user models.User) bool {
	database := infrastructure.GetDB()
	var favorite models.Like
	database.Where(models.Like{
		ArticleId: article.ID,
		UserId:    user.ID,
	}).First(&favorite)
	return favorite.ID != 0
}

func LikeArticle(article *models.Article, user models.User) error {
	database := infrastructure.GetDB()
	var articleLike models.Like
	err := database.FirstOrCreate(&articleLike, &models.Like{
		ArticleId: article.ID,
		UserId:    user.ID,
	}).Error
	return err
}

func Unlike(article *models.Article, user models.User) error {
	database := infrastructure.GetDB()
	err := database.Where(models.Like{
		ArticleId: article.ID,
		UserId:    user.ID,
	}).Delete(models.Like{}).Error
	return err
}

func FetchArticleDetails(condition interface{}, optional ...bool) (models.Article, error) {
	database := infrastructure.GetDB()
	var article models.Article
	tx := database.Begin()

	database.Where(condition).Preload("Tags").Preload("Categories").
		Preload("Comments"). // .Preload("Comments.User")
		First(&article)
	includeUserComment := false

	if len(optional) > 0 {
		includeUserComment = optional[0]
	}

	if includeUserComment {
		// This is not very good from a performance point of view, to improve this we should use select where IN,
		// see my Ecommerce Api app on Github to view the example(on services/products.go#FetchProductDetails)
		for index := range article.Comments {
			tx.Model(&models.Comment{}).Related(&article.Comments[index].User, "User")
		}
	}

	err := tx.Commit().Error
	return article, err
}

func GetComments(article *models.Article) error {
	db := infrastructure.GetDB()
	tx := db.Begin()
	tx.Model(article).Related(&article.Comments, "Comments")
	for i, _ := range article.Comments {
		tx.Model(&article.Comments[i]).Related(&article.Comments[i].User, "User")
		tx.Model(&article.Comments[i].User).Related(&article.Comments[i].User)
	}
	err := tx.Commit().Error
	return err
}

func FetchArticlesPage(page int, pageSize int) ([]models.Article, int, error) {
	db := infrastructure.GetDB()
	var articles []models.Article
	var count int
	tx := db.Begin()
	db.Model(&articles).Count(&count)
	db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&articles)
	tx.Model(&articles).
		Preload("Tags").Preload("Categories").Preload("User").
		Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at desc").
		Find(&articles)

	// TODO: Improve the performance by selecting comment.id where article_id in articleIds
	for i := 0; i < len(articles); i++ {
		articles[i].CommentsCount = tx.Model(&articles[i]).Association("Comments").Count()
	}

	err := tx.Commit().Error
	return articles, count, err
}

func SetTags(model *models.Article, tags []string) error {
	db := infrastructure.GetDB()
	var tagList []models.Tag
	for _, tag := range tags {
		var tagModel models.Tag
		err := db.FirstOrCreate(&tagModel, models.Tag{Name: tag}).Error
		if err != nil {
			return err
		}
		tagList = append(tagList, tagModel)
	}
	model.Tags = tagList
	return nil
}

func UpdateArticle(article *models.Article, data interface{}) error {
	db := infrastructure.GetDB()
	err := db.Model(article).Update(data).Error
	return err
}
func DeleteArticle(condition interface{}) error {
	db := infrastructure.GetDB()
	err := db.Where(condition).Delete(models.Article{}).Error
	return err
}
func DeleteArticleIfOwnerOrAdmin(user *models.User, condition interface{}) error {
	database := infrastructure.GetDB()
	var article models.Article
	err := database.Where(condition).Select("user_id").First(&article).Error
	if err != nil {
		return err
	}
	if user.IsAdmin() || user.ID == article.UserId {
		err = database.Delete(&article).Error
	}
	return err
}

func FetchArticleId(slug string) (uint, error) {
	articleId := -1
	database := infrastructure.GetDB()
	err := database.Model(&models.Article{}).Where(&models.Article{Slug: slug}).Select("id").Row().Scan(&articleId)
	return uint(articleId), err
}
