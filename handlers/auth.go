package handlers

import (
	"bloom/structs"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var week int = 60 * 60 * 24 * 7

func (e *Handlers) LoginHandler(c *gin.Context) {
	var loginCreds structs.LoginJSONPayload
	new := false
	if err := c.ShouldBind(&loginCreds); err != nil {
		MalformedErr(c)
		return
	}

	var user structs.User
	if err := e.DbConn.Where(&structs.User{Email: loginCreds.Email}).First(&user).Error; err != nil {
		// hash the password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(loginCreds.Password), 8)
		confirmID := uuid.NewString()
		// create a new user
		user = structs.User{Email: loginCreds.Email, Password: string(hashedPassword), ConfirmID: confirmID}
		// save
		e.DbConn.Create(&user)
		new = true
		// TK send email to /
		// if the user is found, compare the passwords
	} else if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCreds.Password)); err != nil {
		c.JSON(400, gin.H{"status": "failed", "message": "bad credentials"})
		return
	}
	session := uuid.NewString()

	// TODO make this a transaction
	if err := e.RedisConn.Set("session:"+session, fmt.Sprint(user.ID), 7*24*time.Hour).Err(); err != nil {
		panic(err)
	}
	if err := e.RedisConn.SAdd("sessionsForUser:"+fmt.Sprint(user.ID), session).Err(); err != nil {
		panic(err)
	}

	c.SetCookie("session", session, week, "/", os.Getenv("DOMAIN"), false, true)
	msg := "logged in"
	if new {
		msg = fmt.Sprintf("logged in. You can confirm you email %s at /confirm/%s \n", user.Email, user.ConfirmID)
	}

	c.JSON(200, gin.H{"status": "succeeded", "message": msg})
}

func (e *Handlers) LogoutHandler(c *gin.Context) {
	userId := c.MustGet("userIdAsStr").(string)
	sessions, _ := e.RedisConn.SMembers("sessionsForUser:" + string(userId)).Result()

	for _, sess := range sessions {
		e.RedisConn.Del("session:" + string(sess))
	}
	e.RedisConn.Del("sessionsForUser:" + string(userId))

	c.SetCookie("session", "", 0, "/", "localhost", false, true)
	c.JSON(200, gin.H{"status": "succeeded", "message": "logged out"})
}

func (e *Handlers) SessionMiddleware(c *gin.Context) {
	session, err := c.Cookie("session")
	if err != nil || session == "" {
		c.JSON(400, gin.H{"status": "failed", "message": "not logged in"})
		c.Abort()
		return
	}
	id, err := e.RedisConn.Get("session:" + session).Result()
	if err != nil {
		c.JSON(400, gin.H{"status": "failed", "message": "session not found"})
		c.Abort()
		return
	}

	idAsUint, err := strconv.ParseUint(id, 10, 4)
	if err != nil {
		c.JSON(400, gin.H{"status": "failed", "message": "invalid user id"})
		c.Abort()
		return
	}
	c.Set("userId", idAsUint)
	c.Set("userIdAsStr", id)
	c.Next()
}

func (e *Handlers) WhoAmIHandler(c *gin.Context){
	id := c.MustGet("userId").(uint64)
	var user structs.User
	e.DbConn.Find(&user, id)
	c.JSON(200, gin.H{"email": user.Email, "id": user.ID, "confirmed": user.Confirmed})
}