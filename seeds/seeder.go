package seeds

import (
	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"github.com/melardev/api_blog_app/models"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func randomInt(min, max int) int {

	return rand.Intn(max-min) + min
}

func seedUserRelations(db *gorm.DB) {
	user := models.User{}
	roles := []models.Role{}
	db.Model(&user).Related(&roles, "Roles")
	db.Preload("Roles").First(&user)
	db.Preload("Following").First(&user, "id = ?", 1)

	db.Model(&user).Association("Following").Append(&models.User{FirstName: "friend1"}, &models.User{FirstName: "friend2"})

	db.Model(&user).Association("Following").Delete(&models.User{FirstName: "friend2"})

	db.Model(&user).Association("Following").Replace(&models.User{FirstName: "new friend"})

	db.Model(&user).Association("Following").Clear()

	db.Model(&user).Association("Following").Count()
}
func seedComments(db *gorm.DB) {
	commentsCount := 0
	commentsToSeed := 20

	allUsers := []models.User{}
	allArticles := []models.Article{}
	db.Find(&allArticles)
	db.Find(&allUsers)

	db.Model(&models.Comment{}).Count(&commentsCount)
	commentsToSeed -= commentsCount

	if commentsToSeed > 0 {
		rand.Seed(time.Now().Unix())

		db.Find(&allArticles)
		db.Find(&allUsers)
		for i := 0; i < commentsToSeed; i++ {
			user := allUsers[rand.Intn(len(allUsers))]
			article := allArticles[rand.Intn(len(allArticles))]
			comment := models.Comment{Content: fake.Sentences(), User: user, Article: article}
			db.Set("gorm:association_autoupdate", false).Create(&comment)
		}
	}

	repliesToSeed := 20
	allComments := []models.Comment{}
	allReplies := []models.Comment{}

	// db.Find(&allComments).Preload("Article") // This does not load Article association
	db.Preload("Article").Find(&allComments)
	// for nested eager loading use loading
	// db.Preload("Article").Preload("Article.User").Find(&allComments)
	db.Preload("Article").Where("replied_comment_id IS NOT ?", nil).Find(&allReplies)
	repliesToSeed -= len(allReplies)

	if repliesToSeed > 0 {
		rand.Seed(time.Now().Unix())

		db.Find(&allArticles)
		db.Find(&allUsers)
		for i := 0; i < repliesToSeed; i++ {
			user := allUsers[rand.Intn(len(allUsers))]
			repliedComment := allComments[rand.Intn(len(allComments))]
			comment := models.Comment{Content: fake.Sentences(), User: user, Article: repliedComment.Article,
				RepliedComment: &repliedComment,
			}
			db.Set("gorm:association_autoupdate", false).Create(&comment)
		}
	}

}
func seedAdminFeature(db *gorm.DB) {

	count := 0
	adminRole := models.Role{Name: "ROLE_ADMIN", Description: "Only for admin"}
	q := db.Model(&models.Role{}).Where("name = ?", "ROLE_ADMIN")
	q.Count(&count)

	if count == 0 {
		db.Create(&adminRole)
	} else {
		q.First(&adminRole)
	}

	count = 0
	var adminUsers []models.User
	db.Model(&adminRole).Related(&adminUsers, "Users")

	db.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	if len(adminUsers) == 0 {

		password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		// Approach 1
		//user := models.User{FirstName: "AdminFN", LastName: "AdminFN", Email: "admin@golang.com", Username: "admin", Password: string(password),
		//	Roles: []models.Role{adminRole}}
		//// Do not try to update the adminRole
		////db.Set("gorm:association_autoupdate", false).Create(&user)

		// Approach 2
		user := models.User{FirstName: "AdminFN", LastName: "AdminFN", Email: "admin@golang.com", Username: "admin", Password: string(password)}
		user.Roles = append(user.Roles, adminRole)
		db.NewRecord(user)
		db.Set("gorm:association_autoupdate", false).Save(&user)

		if db.Error != nil {
			print(db.Error)
		}
	}

}

func seedAuthorFeature(db *gorm.DB) {

	count := 0
	authorRole := models.Role{Name: "ROLE_AUTHOR", Description: "Only for authors"}
	q := db.Model(&models.Role{}).Where("name = ?", "ROLE_AUTHOR")
	q.Count(&count)

	if count == 0 {
		db.Create(&authorRole)
	} else {
		q.First(&authorRole)
	}

	count = 0
	var authors []models.User

	db.Model(&authorRole).Related(&authors, "Users")

	authorsCount := len(authors)
	authorsToSeed := 5
	authorsToSeed -= authorsCount
	if authorsToSeed > 0 {
		password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		for i := 0; i < authorsToSeed; i++ {
			user := models.User{FirstName: fake.FirstName(), LastName: fake.LastName(), Email: fake.EmailAddress(), Username: fake.UserName(),
				Password: string(password)}
			// No need to add the role as we did for seedAdmin, it is added by the BeforeSave hook
			user.Roles = append(user.Roles, authorRole)
			db.NewRecord(user)
			db.Set("gorm:association_autoupdate", false).Save(&user)
			//	db.Create(&user)

			if db.Error != nil {
				print(db.Error)
			}
		}
	}
}

func seedUsers(db *gorm.DB) {

	count := 0
	role := models.Role{Name: "ROLE_USER", Description: "Only for standard users"}
	q := db.Model(&models.Role{}).Where("name = ?", "ROLE_USER")
	q.Count(&count)

	if count == 0 {
		db.Create(&role)
	} else {
		q.First(&role)
	}

	var standardUsers []models.User
	db.Model(&role).Related(&standardUsers, "Users")
	usersCount := len(standardUsers)
	usersToSeed := 20
	usersToSeed -= usersCount
	if usersToSeed > 0 {
		password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		for i := 0; i < usersToSeed; i++ {
			user := models.User{FirstName: fake.FirstName(), LastName: fake.LastName(), Email: fake.EmailAddress(), Username: fake.UserName(),
				Password: string(password)}
			// No need to add the role as we did for seedAdmin, it is added by the BeforeSave hook
			db.Create(&user)
		}
	}

}
func seedTags(db *gorm.DB) {
	db.Where(&models.Tag{Name: "Spring MVC"}).Attrs(models.Tag{Description: "Articles for Spring MVC", IsNewRecord: true}).FirstOrCreate(&models.Tag{})
	db.Where(models.Tag{Name: "Laravel"}).Attrs(models.Tag{Description: "Articles for Laravel", IsNewRecord: true}).FirstOrCreate(&models.Tag{})
	db.Where(models.Tag{Name: "Ruby on Rails"}).Attrs(models.Tag{Description: "Articles of Ruby on Rails", IsNewRecord: true}).FirstOrCreate(&models.Tag{})
}

func seedCategories(db *gorm.DB) {
	db.Where(models.Category{Name: "Cpp"}).Attrs(models.Category{Description: "Articles for Cpp", IsNewRecord: true}).FirstOrCreate(&models.Category{})
	db.Where(models.Category{Name: "Java"}).Attrs(models.Category{Description: "Articles for Java", IsNewRecord: true}).FirstOrCreate(&models.Category{})
	db.Where(models.Category{Name: "Ruby"}).Attrs(models.Category{Description: "Articles of Ruby", IsNewRecord: true}).FirstOrCreate(&models.Category{})
}

func seedArticles(db *gorm.DB) {

	authorRole := models.Role{}
	db.Model(&models.Role{}).Where("name = ?", "ROLE_AUTHOR").First(&authorRole)

	var authors []models.User
	db.Model(&authorRole).Related(&authors, "Users")

	articlesCount := 0
	articlesToSeed := 20
	db.Model(&models.Article{}).Count(&articlesCount)
	articlesToSeed -= articlesCount

	if articlesToSeed > 0 {
		rand.Seed(time.Now().Unix())
		tags := []models.Tag{}
		categories := []models.Category{}
		db.Find(&tags)
		db.Find(&categories)
		for i := 0; i < articlesToSeed; i++ {
			tagsForArticle := tags[rand.Intn(len(tags))]
			categoriesForArticle := categories[rand.Intn(len(tags))]

			article := &models.Article{Title: fake.WordsN(randomInt(2, 4)), Description: fake.Paragraph(), Body: fake.Paragraphs(),
				User: authors[rand.Intn(len(authors))], Tags: []models.Tag{tagsForArticle},
				Categories: []models.Category{categoriesForArticle}}

			db.Set("gorm:association_autoupdate", false).Create(&article)
		}
	}
}

func seedUserSubscriptions(db *gorm.DB) {

	var authors []models.User
	var allUsers []models.User
	var subscriptions []models.UserSubscription

	db.Find(&allUsers)
	db.Find(&subscriptions)
	db.Table("users_subscriptions").Find(&subscriptions)

	for i := 0; i < len(allUsers); i++ {
		user := allUsers[i]
		for _, role := range user.Roles {
			if role.Name == "ROLE_AUTHOR" {
				authors = append(authors, user)
			}
		}
	}

	subscriptionsCount := len(subscriptions)
	subscriptionsToSeed := 10
	subscriptionsToSeed -= subscriptionsCount
	if subscriptionsToSeed > 0 {
		rand.Seed(time.Now().Unix())

		for i := 0; i < subscriptionsToSeed; i++ {
			following := allUsers[rand.Intn(len(allUsers))]
			follower := allUsers[rand.Intn(len(allUsers))]
			subscription := models.UserSubscription{Following: following, Follower: follower}
			db.Set("gorm:association_autoupdate", false).Create(&subscription)
		}
	}
}

func seedSiteSubscriptions(db *gorm.DB) {

	var subscriptions []models.SiteSubscription
	var usersNotSubscribedYet []models.User
	var ids []int64
	db.Select("user_id").Find(&subscriptions).Pluck("user_id", &ids)
	// db.Select("id").Not("id", ids).Find(&usersNotSubscribedYet)
	db.Select("id").Find(&usersNotSubscribedYet)

	subscriptionsCount := len(subscriptions)
	subscriptionsToSeed := 5
	subscriptionsToSeed -= subscriptionsCount
	if subscriptionsToSeed > 0 {
		rand.Seed(time.Now().Unix())

		for i := 0; i < subscriptionsToSeed; i++ {
			user := usersNotSubscribedYet[rand.Intn(len(usersNotSubscribedYet))]
			subscription := models.SiteSubscription{User: user}
			db.Set("gorm:association_autoupdate", false).Create(&subscription)
		}
	}

}
func seedLikes(db *gorm.DB) {
	likesCount := 0
	linksToSeed := 20

	allUsers := []models.User{}
	allLikes := []models.Article{}
	db.Find(&allLikes)
	db.Find(&allUsers)

	db.Model(&models.Like{}).Count(&likesCount)
	linksToSeed -= likesCount

	if linksToSeed > 0 {
		rand.Seed(time.Now().Unix())

		db.Find(&allLikes)
		db.Find(&allUsers)
		for i := 0; i < linksToSeed; i++ {
			user := allUsers[rand.Intn(len(allUsers))]
			article := allLikes[rand.Intn(len(allLikes))]
			like := models.Like{User: user, Article: article}
			db.Set("gorm:association_autoupdate", false).Create(&like)
		}
	}
}

func Seed(db *gorm.DB) {
	var user models.User
	var followers []models.User
	db.First(&user)
	db.Model(&user).Related(&followers, "Following")

	seedAdminFeature(db)
	seedAuthorFeature(db)
	seedUsers(db)
	seedTags(db)
	seedCategories(db)
	seedArticles(db)
	seedComments(db)
	seedLikes(db)
	seedUserSubscriptions(db)
	seedSiteSubscriptions(db)
}
