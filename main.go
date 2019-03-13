package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/melardev/GoGonicBlogApi/controllers"
	"github.com/melardev/GoGonicBlogApi/infrastructure"
	"github.com/melardev/GoGonicBlogApi/middlewares"
	"github.com/melardev/GoGonicBlogApi/models"
	"github.com/melardev/GoGonicBlogApi/seeds"
	"os"
)

func drop(db *gorm.DB) {
	db.DropTableIfExists(&models.Like{},
		models.ArticleTag{}.TableName(), models.ArticleCategory{}.TableName(),
		&models.Tag{}, &models.Category{},
		&models.Comment{}, &models.Article{},
		&models.UserSubscription{}, &models.SiteSubscription{},
		&models.UserRole{}, &models.Role{}, &models.User{})
}
func migrate(db *gorm.DB) {

	db.AutoMigrate(&models.UserSubscription{}) // UserSubscription must be migrated first otherwise the join table create has not the shape we are expecting
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Article{})
	db.AutoMigrate(&models.Tag{})
	db.AutoMigrate(&models.Like{})
	db.AutoMigrate(&models.Comment{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Role{})

	db.AutoMigrate(&models.SiteSubscription{})

}

func addDbConstraints(database *gorm.DB) {
	// TODO: it is well known GORM does not add foreign keys even after using ForeignKey in struct, but, why manually does not work neither ??

	dialect := database.Dialect().GetName() // mysql, sqlite3
	if dialect != "sqlite3" {
		database.Model(&models.Comment{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
		database.Model(&models.Comment{}).AddForeignKey("article_id", "articles(id)", "CASCADE", "CASCADE")

		database.Model(&models.Article{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
		database.Model(&models.Article{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

		database.Model(&models.Comment{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

		database.Model(&models.Like{}).AddForeignKey("article_id", "articles(id)", "CASCADE", "CASCADE")
		database.Model(&models.Like{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

		database.Model(&models.UserRole{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
		database.Model(&models.UserRole{}).AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")

		database.Model(&models.ArticleTag{}).AddForeignKey("article_id", "articles(id)", "CASCADE", "CASCADE")
		// We can pass the table name string to Model as well, but to Table() we can only pass the table name string, and not the Model reference
		database.Model(models.ArticleTag{}.TableName()).AddForeignKey("tag_id", "tags(id)", "CASCADE", "CASCADE")

		database.Table(models.ArticleCategory{}.TableName()).AddForeignKey("article_id", "articles(id)", "CASCADE", "CASCADE")
		database.Model(&models.ArticleCategory{}).AddForeignKey("category_id", "categories(id)", "CASCADE", "CASCADE")

		database.Model(&models.UserSubscription{}).AddForeignKey("following_id", "users(id)", "CASCADE", "CASCADE")
		database.Table("users_subscriptions").AddForeignKey("follower_id", "users(id)", "CASCADE", "CASCADE")

	} else if dialect == "sqlite3" {
		database.Table("comments").AddIndex("comments__idx_article_id", "article_id")
		database.Table("comments").AddIndex("comments__idx_user_id", "user_id")

		database.Model(&models.Comment{}).AddIndex("comments__idx_created_at", "created_at")

	}

	database.Model(&models.UserRole{}).AddIndex("user_roles__idx_user_id", "user_id")
	database.Table("articles_tags").AddIndex("articles_tags__idx_article_id", "article_id")
}

func create(database *gorm.DB) {
	drop(database)
	migrate(database)
	addDbConstraints(database)
}

func main() {
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	databse := infrastructure.InitDb()
	defer databse.Close()
	args := os.Args
	if len(args) > 1 {
		first := args[1]
		second := ""
		if len(args) > 2 {
			second = args[2]
		}

		if first == "create" {
			create(databse)
		} else if first == "seed" {
			seeds.Seed(databse)
			os.Exit(0)
		} else if first == "migrate" {
			migrate(databse)
		}

		if second == "seed" {
			seeds.Seed(databse)
			os.Exit(0)
		} else if first == "migrate" {
			migrate(databse)
		}

		if first != "" && second == "" {
			os.Exit(0)
		}
	}
	migrate(databse)
	router := gin.Default() // gin with the Logger and Recovery Middlewares attached
	router.Use(middlewares.Benchmark())
	router.Use(cors.Default())
	router.Use(middlewares.UserLoaderMiddleware())
	apiRouteGroup := router.Group("/api")

	controllers.RegisterUserRoutes(apiRouteGroup.Group("/users"))
	controllers.RegisterArticleRoutes(apiRouteGroup.Group("/articles"))
	controllers.RegisterCommentRoutes(apiRouteGroup.Group("/"))
	controllers.RegisterLikeRoutes(apiRouteGroup.Group("/"))
	controllers.RegisterTagRoutes(apiRouteGroup.Group("/tags"))
	controllers.RegisterCategoryRoutes(apiRouteGroup.Group("/categories"))
	controllers.RegisterUserSubscriptionRoutes(apiRouteGroup.Group("/users"))

	router.Run(":8080") // listen and serve on 0.0.0.0:8080
}
