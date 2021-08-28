package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"bloom/email"
	"bloom/structs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var week int = 60 * 60 * 24 * 7

func (e *Handlers) SignupHandler(c *gin.Context) {
	var loginCreds structs.LoginJSONPayload
	if err := c.ShouldBind(&loginCreds); err != nil {
		MalformedErr(c)
		return
	}

	var user structs.User
	err := e.DbConn.Where(&structs.User{Email: loginCreds.Email}).First(&user).Error

	if err == nil {
		UserAlreadyExistsErr(c)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(loginCreds.Password), 8)
	confirmID := uuid.NewString()
	// create a new user
	user = structs.User{Email: loginCreds.Email, Password: string(hashedPassword), ConfirmID: confirmID}
	// save
	e.DbConn.Create(&user)
	email.SendConfirmEmail(user.Email, user.ConfirmID)

	session := uuid.NewString()
	// TODO make this a transaction
	if err := e.RedisConn.Set("session:"+session, fmt.Sprint(user.ID), 7*24*time.Hour).Err(); err != nil {
		panic(err)
	}
	if err := e.RedisConn.SAdd("sessionsForUser:"+fmt.Sprint(user.ID), session).Err(); err != nil {
		panic(err)
	}
	c.SetCookie("session", session, week, "/", os.Getenv("DOMAIN"), true, true)
	c.JSON(200, structs.JsonResponse{Succeeded: true, Message: "New user has been created. Please check your email for a confirmation."})
}

func (e *Handlers) LoginHandler(c *gin.Context) {
	var loginCreds structs.LoginJSONPayload
	if err := c.ShouldBind(&loginCreds); err != nil {
		MalformedErr(c)
		return
	}

	var user structs.User
	if err := e.DbConn.Where(&structs.User{Email: loginCreds.Email}).First(&user).Error; err != nil {
		UserDoesntExistErr(c)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCreds.Password)); err != nil {
		InvalidCredsErr(c)
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
	c.SetCookie("session", session, week, "/", os.Getenv("DOMAIN"), true, true)
	c.JSON(200, structs.JsonResponse{Succeeded: true})
}

func (e *Handlers) LogoutHandler(c *gin.Context) {
	userId := c.MustGet("userIdAsStr").(string)
	sessions, _ := e.RedisConn.SMembers("sessionsForUser:" + string(userId)).Result()

	for _, sess := range sessions {
		e.RedisConn.Del("session:" + string(sess))
	}
	e.RedisConn.Del("sessionsForUser:" + string(userId))
	c.SetCookie("session", "", 0, "/", os.Getenv("DOMAIN"), true, true)
	// https://github.com/gin-gonic/gin#redirects
	c.Redirect(http.StatusFound, os.Getenv("FRONTEND_URL")+"?msg="+url.QueryEscape("You have been logged out"))
}

func (e *Handlers) LoginHandlerOld(c *gin.Context) {
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
		email.SendConfirmEmail(user.Email, user.ConfirmID)
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
	c.SetCookie("session", session, week, "/", os.Getenv("DOMAIN"), true, true)
	msg := "logged in"
	if new {
		msg = fmt.Sprintf("signed up. You can confirm you email %s at /confirm/%s \n", user.Email, user.ConfirmID)
	}

	c.JSON(200, gin.H{"status": "succeeded", "message": msg, "new": new})
}

func (e *Handlers) SessionMiddleware(c *gin.Context) {
	session, err := c.Cookie("session")
	if err != nil || session == "" {
		NotLoggedInErr(c)
		c.Abort()
		return
	}
	id, err := e.RedisConn.Get("session:" + session).Result()
	if err != nil {
		SessionNotFoundErr(c)
		c.Abort()
		return
	}

	idAsUint, err := strconv.ParseUint(id, 10, 4)
	c.Set("userId", idAsUint)
	c.Set("userIdAsStr", id)
	c.Next()
}

func (e *Handlers) WhoAmIHandler(c *gin.Context) {
	id := c.MustGet("userId").(uint64)
	var user structs.User
	e.DbConn.Find(&user, id)
	userData := gin.H{"email": user.Email, "id": user.ID, "confirmed": user.Confirmed}
	c.JSON(200, structs.JsonResponse{Succeeded: true, Data: userData})
}
