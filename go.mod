module bloom

// +heroku goVersion go1.116
go 1.16

require (
	github.com/gin-gonic/gin v1.7.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.6.0 // indirect
	golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.13
)
