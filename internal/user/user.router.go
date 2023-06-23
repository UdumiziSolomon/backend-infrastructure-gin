package user

import (
	"github.com/gin-gonic/gin"
)


func UserRoutes(router *gin.Engine){

	user_router := router.Group("/api/user")

	user_router.GET("/users", GetUsers())
	user_router.POST("/new-user", CreateUser())

}