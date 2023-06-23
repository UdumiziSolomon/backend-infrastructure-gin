package main

import (
	"log"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"github.com/solonode/golang-jwt-mongo/internal/database"
	"github.com/solonode/golang-jwt-mongo/internal/user"
	"github.com/solonode/golang-jwt-mongo/config"
)

func main(){

	// Loading env vars
	config, err := config.LoadENVFile("../../")  // file path of app.env
	if err != nil{
		log.Fatal("Error loading env vars: ", err)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// cors config
	corsConfig := cors.DefaultConfig()
	// corsConfig.AllowAllOrigins = true // allow all origins(*)
	corsConfig.AllowOrigins = []string{config.ClientUrl} // allow from client url
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	

	router.Use(cors.New(corsConfig))   // bind cors config to new cor instance

	// MongoDB connection
	client := database.SetupMongoDB()

	defer client.Disconnect(context.Background())  // close the DB connection


	// Routes declr
	user.UserRoutes(router)

	router.Run(config.ServerPort)  // server initialization

}	