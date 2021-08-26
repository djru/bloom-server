package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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
// TODO find a way to keep the heroku app up
// TODO set up stack: vercel, heroku, sendgrid, cloudflare
// TODO endpoint to get all readings as csv
// TODO endpoint to delete my data (send email)

var dbConn *gorm.DB
var ctx = context.Background()
var redisConn *redis.Client
var week int = 60 * 60 * 24 * 7

func setupHandlers(r *gin.Engine, db *gorm.DB, redis *redis.Client) {
	h := handlers.Handlers{DbConn: db, RedisConn: redis}
	r.GET("/", h.HomeHandler)

	r.GET("/resend", h.SessionMiddleware, h.ReSendConfirmEmailHandler)
	r.GET("/confirm/:id", h.ConfirmEmailHandler)
	r.GET("/sendRecover", h.StartRecoveryProcessHandler)
	r.POST("/recover", h.EndRecoveryProcessHandler)

	r.GET("/readings", h.SessionMiddleware, h.GetReadingsHandler)
	r.POST("/newReading", h.SessionMiddleware, h.NewReadingHandler)

	r.POST("/login", h.LoginHandler)
	r.POST("/signup", h.LoginHandler)
	r.GET("/logout", h.SessionMiddleware, h.LogoutHandler)

	r.GET("/whoami", h.SessionMiddleware, h.WhoAmIHandler)
}

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())

	config := cors.DefaultConfig()
	// https://github.com/gin-contrib/cors#using-defaultconfig-as-start-point
	// https://stackoverflow.com/questions/29418478/go-gin-framework-cors
	// https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
	// https://stackoverflow.com/a/48763475/5360657
	// https://blog.heroku.com/chrome-changes-samesite-cookie

	// A brief note on how the app is being hosted in prod.
	// Heroku makes you pay for ssl support, which I don't want to do. But cloudflare can proxy
	// the requests, do ssl termination and then do http to the backend
	// this is ideal for the heroku setup BUT it breaks the vercel setup.
	// This is because vercel automatically redirects http to https
	// which causes a redirect loop because ever time the browser is redirected, cloudflare turns it into an http request
	// so I had to turn off cloudflare proxing for the frontend.
	config.AllowOrigins = []string{"https://bloom-health.herokuapp.com", "https://bloom-ui.vercel.app", "https://www.bloomhealth.app", "https://bloomhealth.app"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	if os.Getenv("ENV") == "dev" {
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

	// https://pkg.go.dev/github.com/go-redis/redis#ParseURL
	fmt.Println(os.Getenv("REDIS_URL"))
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
	// https://stackoverflow.com/questions/28706180/setting-the-port-for-node-js-server-on-heroku
	log.Fatal(r.Run(":" + os.Getenv("PORT")))
}
