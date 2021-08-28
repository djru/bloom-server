package handlers

import (
	"bloom/email"
	"bloom/structs"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TODO make this post using json
func (e *Handlers) StartRecoveryProcessHandler(c *gin.Context) {
	em := c.DefaultQuery("email", "")
	if em == "" {
		InvalidCredsErr(c)
		return
	}

	var user structs.User
	err := e.DbConn.First(&user, "email = ?", em).Error
	if err != nil {
		UserDoesntExistErr(c)
		return
	}
	user.RecoveryID = uuid.NewString()
	// set the expiry to be 20 minutes in the future
	user.RecoveryExpiry = time.Now().Add(time.Minute * 20).Unix()
	// TK send email
	e.DbConn.Save(&user)
	email.SendRecoveryEmail(user.Email, user.RecoveryID)
	c.JSON(http.StatusAccepted, structs.JsonResponse{Succeeded: true, Message: "A recovery link has been sent to " + user.Email})
}

func (e *Handlers) EndRecoveryProcessHandler(c *gin.Context) {
	var payload structs.RecoveryJSONPayload
	var user structs.User
	if err := c.ShouldBind(&payload); err != nil {
		InvalidCredsErr(c)
		return
	}
	if err := e.DbConn.Where(&structs.User{RecoveryID: payload.RecoveryID}).First(&user).Error; err != nil {
		SessionNotFoundErr(c)
		return
	}
	if user.Email != payload.Email {
		UserDoesntExistErr(c)
		return
	}

	if user.RecoveryExpiry < time.Now().Unix() {
		ExpiredErr(c)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), 8)
	user.Password = string(hashedPassword)
	e.DbConn.Save(&user)
	c.JSON(200, structs.JsonResponse{Succeeded: true, Message: "Password reset. Please log in again"})
}

func (e *Handlers) ConfirmEmailHandler(c *gin.Context) {
	var user structs.User
	id := c.Param("id")
	err := e.DbConn.Where(&structs.User{ConfirmID: id}).First(&user).Error
	if err != nil {
		c.JSON(401, gin.H{"message": "no user found"})
		return
	}
	user.ConfirmID = ""
	user.Confirmed = true
	e.DbConn.Save(&user)
	c.Redirect(http.StatusFound, os.Getenv("FRONTEND_URL")+"?msg="+url.QueryEscape("Email confirmed"))
}

func (e *Handlers) ReSendConfirmEmailHandler(c *gin.Context) {
	id := c.MustGet("userId").(uint64)
	var user structs.User
	e.DbConn.Find(&user, id)
	if user.Confirmed {
		c.Redirect(http.StatusFound, os.Getenv("FRONTEND_URL")+"/me?msg="+url.QueryEscape("Your email is already confirmed"))
		return
	}
	email.SendConfirmEmail(user.Email, user.ConfirmID)
	c.Redirect(http.StatusFound, os.Getenv("FRONTEND_URL")+"/me?msg="+url.QueryEscape("Confirmation email sent"))
}
