package handlers

import (
	"bloom/email"
	"bloom/structs"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (e *Handlers) StartRecoveryProcessHandler(c *gin.Context) {
	em := c.DefaultQuery("email", "")
	if em == "" {
		c.JSON(400, gin.H{"message": "invalad params"})
		return
	}

	var user structs.User
	err := e.DbConn.First(&user, "email = ?", em).Error
	if err != nil {
		c.JSON(401, gin.H{"message": "that email is not associated with any user"})
		return
	}
	user.RecoveryID = uuid.NewString()

	// TK send email
	e.DbConn.Save(&user)
	c.JSON(200, gin.H{"message": fmt.Sprintf("You can recover password for %s at /recover with recovery_id:'%v' \n", user.Email, user.RecoveryID)})
	email.SendRecoveryEmail(user.Email, user.RecoveryID)

}

func (e *Handlers) EndRecoveryProcessHandler(c *gin.Context) {
	var payload structs.RecoveryJSONPayload
	var user structs.User
	if err := c.ShouldBind(&payload); err != nil {
		c.JSON(400, gin.H{"message": "invalid payload"})
		return
	}
	if err := e.DbConn.Where(&structs.User{RecoveryID: payload.RecoveryID}).First(&user).Error; err != nil {
		c.JSON(400, gin.H{"message": "no user found"})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), 8)
	user.Password = string(hashedPassword)
	e.DbConn.Save(&user)
	c.JSON(200, gin.H{"message": "new password saved"})
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
	c.JSON(200, gin.H{"message": "user confirmed"})
}

func (e *Handlers) ReSendConfirmEmail(c *gin.Context) {
	id := c.MustGet("userId").(uint64)
	var user structs.User
	e.DbConn.Find(&user, id)
	email.SendConfirmEmail(user.Email, user.ConfirmID)
	c.JSON(200, gin.H{"email": user.Email, "id": user.ID, "confirmed": user.Confirmed})
}
