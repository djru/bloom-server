package structs

import "gorm.io/gorm"

type LoginJSONPayload struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RecoveryJSONPayload struct {
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
