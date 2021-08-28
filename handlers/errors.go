package handlers

import (
	"bloom/structs"

	"github.com/gin-gonic/gin"
)

func MalformedErr(c *gin.Context) {
	c.JSON(400, structs.JsonResponse{Succeeded: false, Message: "Bad payload"})
	return
}

func UserAlreadyExistsErr(c *gin.Context) {
	c.JSON(400, structs.JsonResponse{Succeeded: false, Message: "A user with that email already exists"})
	return
}

func UserDoesntExistErr(c *gin.Context) {
	c.JSON(400, structs.JsonResponse{Succeeded: false, Message: "A user with that email doesn't exist"})
	return
}

func InvalidCredsErr(c *gin.Context) {
	c.JSON(400, structs.JsonResponse{Succeeded: false, Message: "The password does not match"})
	return
}

func SessionNotFoundErr(c *gin.Context) {
	c.JSON(400, structs.JsonResponse{Succeeded: false, Message: "The session was not found"})
	return
}

func NotLoggedInErr(c *gin.Context) {
	c.JSON(400, structs.JsonResponse{Succeeded: false, Message: "You are not logged in"})
	return
}
