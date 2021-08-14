package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"bloom/handlers"
	"bloom/structs"
)

// TODO email verification on signup
// TODO domain name
// TODO deploy the front end to vercel
// TODO deploy the back end to heroku
// TODO find a way to keep the heroku app up
// TODO CORS if necessary
// TODO set up stack: vercel, heroku, sendgrid, cloudflare
// TODO endpoint to get all readings as csv
// TODO endpoint to delete my data (send email)
// TODO extract claims in a function, abstract it and use it anywhere we need

var identityKey = "email"
var dbConn *gorm.DB
var ctx = context.Background()
var redisConn *redis.Client
var week int = 60 * 60 * 24 * 7


func setupHandlers(r *gin.Engine, db *gorm.DB, redis *redis.Client){
	h := handlers.Handlers{DbConn: db, RedisConn: redis}
	r.GET("/", h.HomeHandler)

	r.GET("/confirm/:id", h.ConfirmEmailHandler)
	r.GET("/sendRecover", h.StartRecoveryProcessHandler)
	r.POST("/recover", h.EndRecoveryProcessHandler)

	r.GET("/readings", h.GetReadingsHandler)
	r.POST("/newReading", h.NewReadingHandler)

	r.POST("/login", h.LoginHandler)
	r.GET("/logout", h.SessionMiddleware, h.LogoutHandler)

	r.GET("/whoami", h.SessionMiddleware, h.WhoAmIHandler)
}

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())

	if os.Getenv("ENV") == "dev"{
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	db, _ := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	sqlDB, _ := db.DB()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	db.AutoMigrate(&structs.User{})
	db.AutoMigrate(&structs.Reading{})

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	fmt.Println("addr is", opt.Addr)
	fmt.Println("db is", opt.DB)
	fmt.Println("password is", opt.Password)

	dbConn = db
	redisConn = redis.NewClient(opt)


	setupHandlers(r, dbConn, redisConn)
	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}
