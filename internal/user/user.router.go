package user

import (
	"github.com/gin-gonic/gin"
)


func UserRoutes(router *gin.Engine){

	user_router := router.Group("/api/user")

	user_router.POST("/verify-user", VerifyUser())
	user_router.POST("/new-user", CreateUser())
	user_router.POST("/login", LoginUser())
	user_router.GET("/user/:uid", GetSingleUser())
	user_router.GET("/users", GetAllUsers())
	user_router.PUT("/update/:uid", UpdateSingleUser())
	user_router.DELETE("/delete/:uid", DeleteSingleUser())
	user_router.DELETE("/delete/users", DeleteAllUsers())

}