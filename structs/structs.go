package structs

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginJSONPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RecoveryJSONPayload struct {
	Email      string `json:"email" binding:"required"`
	RecoveryID string `json:"recovery_id" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type NewReadingJSONPayload struct {
	Top     int `json:"top" binding:"required"`
	Bottom  int `json:"bottom" binding:"required"`
	Pulse   int `json:"pulse"`
	Feeling int `json:"feeling"`
}

type User struct {
	gorm.Model
	Email      string `gorm:"index"`
	Password   string
	Confirmed  bool `gorm:"default:false"`
	ConfirmID  string
	RecoveryID string
}

func (user *User) GetReturnableData() gin.H {
	return gin.H{"email": user.Email, "id": user.ID, "confirmed": user.Confirmed}
}

type Reading struct {
	gorm.Model
	Top     int
	Bottom  int
	Pulse   int
	Feeling int
	UserID  uint
	User    User `json:"-"`
}

type JsonResponse struct {
	Succeeded bool                   `json:"succeeded"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data"`
}
