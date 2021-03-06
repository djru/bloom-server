package handlers

import (
	"bloom/structs"

	"github.com/gin-gonic/gin"
)

func (e *Handlers) NewReadingHandler(c *gin.Context) {
	id := c.MustGet("userId").(uint64)
	var payload structs.NewReadingJSONPayload
	if err := c.ShouldBind(&payload); err != nil {
		MalformedErr(c)
		return
	}

	reading := structs.Reading{Top: payload.Top, Bottom: payload.Bottom, Pulse: payload.Pulse, UserID: uint(id), Feeling: payload.Feeling}
	e.DbConn.Create(&reading)

	c.JSON(200, structs.JsonResponse{Succeeded: true, Data: reading})
}

func (e *Handlers) GetReadingsHandler(c *gin.Context) {
	id := c.MustGet("userId").(uint64)
	var readings []structs.Reading
	e.DbConn.Where(&structs.Reading{UserID: uint(id)}).Find(&readings)
	c.JSON(200, structs.JsonResponse{Succeeded: true, Data: readings})
}
