module bloom

// https://stackoverflow.com/questions/56968852/specify-go-version-for-go-mod-file
// +heroku goVersion go1.16
go 1.16

require (
	github.com/gin-contrib/cors v1.3.1 // indirect
	github.com/gin-gonic/gin v1.7.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.3.0
	github.com/lib/pq v1.6.0 // indirect
	github.com/sendgrid/rest v2.6.4+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.10.0+incompatible // indirect
	golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.13
)
