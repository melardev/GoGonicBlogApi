package main

import (
	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"github.com/melardev/GoGonicBlogApi/models"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

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

	articlesCount := 0
	articlesToSeed := 20
	db.Model(&models.Article{}).Count(&articlesCount)
	articlesToSeed -= articlesCount

	if articlesCount > 0 {
		rand.Seed(time.Now().Unix())
		allUsers := []models.User{}
		allArticles := []models.Article{}
		db.Find(&allUsers)
		db.Find(&allArticles)
		for i := 0; i < articlesCount; i++ {
			user := allUsers[rand.Intn(len(allUsers))]
			article := allArticles[rand.Intn(len(allArticles))]
			comment := models.Comment{Content: fake.Sentences(), User: user, Article: article}
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

	db.Model(&models.User{}).Where("username = ?", "=", "admin").Count(&count)
	if len(adminUsers) == 0 {
		adminRole := models.Role{}
		q.First(&adminRole)
		password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		// Approach 1
		user := models.User{FirstName: "AdminFN", LastName: "AdminFN", Email: "admin@golang.com", Username: "admin", Password: string(password),
			Roles: []models.Role{adminRole}}
		// Do not try to update the adminRole
		db.Set("gorm:association_autoupdate", false).Create(&user)

		// Approach 2
		// user := users.User{FirstName: "AdminFN", LastName: "AdminFN", Email: "admin@golang.com", Username: "admin", Password: string(password)}
		// user.Roles = append(user.Roles, adminRole)
		// db.NewRecord(user)
		// db.Set("gorm:association_autoupdate", false).Save(&user)

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
		// Approach 2
		for i := 0; i < authorsToSeed; i++ {
			password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
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
		for i := 0; i < usersToSeed; i++ {
			user := models.User{FirstName: fake.FirstName(), LastName: fake.LastName(), Email: fake.EmailAddress(), Username: fake.UserName()}
			// No need to add the role as we did for seedAdmin, it is added by the BeforeSave hook
			db.Create(&user)
		}
	}

}
func seedTags(db *gorm.DB) {

	db.Where(&models.Tag{Name: "Spring MVC"}).Attrs(models.Tag{Description: "Articles for Spring MVC"}).FirstOrCreate(&models.Tag{})
	db.Where(models.Tag{Name: "Laravel"}).Attrs(models.Tag{Description: "Articles for Laravel"}).FirstOrCreate(&models.Tag{})
	db.Where(models.Tag{Name: "Ruby on Rails"}).Attrs(models.Tag{Description: "Articles of Ruby on Rails"}).FirstOrCreate(&models.Tag{})
}

func seedCategories(db *gorm.DB) {
	db.Where(models.Category{Name: "Cpp"}).Attrs(models.Category{Description: "Articles for Cpp"}).FirstOrCreate(&models.Category{})
	db.Where(models.Category{Name: "Java"}).Attrs(models.Category{Description: "Articles for Java"}).FirstOrCreate(&models.Category{})
	db.Where(models.Category{Name: "Ruby"}).Attrs(models.Category{Description: "Articles of Ruby"}).FirstOrCreate(&models.Category{})
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

	if articlesCount > 0 {
		rand.Seed(time.Now().Unix())
		tags := []models.Tag{}
		categories := []models.Category{}
		db.Find(&tags)
		db.Find(&categories)
		for i := 0; i < articlesCount; i++ {
			tagsForArticle := tags[rand.Intn(len(tags))]
			categoriesForArticle := categories[rand.Intn(len(tags))]

			article := &models.Article{Title: fake.Sentence(), Description: fake.Paragraph(), Body: fake.Paragraphs(),
				User: authors[rand.Intn(len(authors))], Tags: []models.Tag{tagsForArticle},
				Categories: []models.Category{categoriesForArticle}}

			db.Set("gorm:association_autoupdate", false).Create(&article)
		}
	}
}

func Seed(db *gorm.DB) {
	seedAdminFeature(db)
	seedAuthorFeature(db)
	seedUsers(db)
	seedTags(db)
	seedCategories(db)
	seedArticles(db)
	seedComments(db)
}
