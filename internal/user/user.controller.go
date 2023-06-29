package user 

import(
	"context"
	"net/http"
	"time"
	"log"
	"fmt"
	"math/rand"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/solonode/golang-jwt-mongo/internal/database"
	"github.com/solonode/golang-jwt-mongo/config"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "user")

func GetSingleUser() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10 * time.Second)
		defer cancel()

		userID := c.Param("uid")
		id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var user User  

		err = userCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": user})
	}
}


func GetAllUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10 * time.Second)
		defer cancel()

		cursor, err := userCollection.Find(ctx, bson.D{{}})
		if err != nil {
			log.Fatal(err)
		}

		var users []primitive.M
		for cursor.Next(ctx){
			var user bson.M
			if err := cursor.Decode(&user); err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}
		defer cursor.Close(ctx)

		c.JSON(http.StatusOK, gin.H{"data": users})

	}
}

func VerifyUser() gin.HandlerFunc {
	return func(c *gin.Context){
		var newUser User   // request body

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10 * time.Second)
		defer cancel()

		// validate the request body
		if err := c.ShouldBindJSON(&newUser); err != nil{
			c.JSON(http.StatusBadRequest, err)
			return
		}

		// check if user exists
		check, err := userCollection.CountDocuments(ctx, bson.M{"email": newUser.Email})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		if check > 0 {
			c.JSON(http.StatusConflict, gin.H{"update": "User already exist!"})
			return
		}

			// Loading env vars
		config, err := config.LoadENVFile("../../")  // file path of app.env
		if err != nil{
			log.Fatal("Error loading env vars: ", err)
		}

		// mail
		rand.Seed(time.Now().UnixNano())
		verificationCode := rand.Intn(900000) + 100000

		toAddr := newUser.Email  // receiver

		// message
		subject := "Account verification code"
		body := fmt.Sprintf("<h2> Verify your account with this code </h2><h1> %v </h1> ", verificationCode)

		m := gomail.NewMessage()  // new instance for gomail
		m.SetHeader("From", config.MailSender)
		m.SetHeader("To", *toAddr)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)

		mail := gomail.NewDialer(config.SmtpHost, config.SmtpPort, config.MailSender, config.MailPassword)
		if err := mail.DialAndSend(m); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Mail was sent"})

	}
}


func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context){
		var newUser User  

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10 * time.Second)
		defer cancel()

		// validate the request body
		if err := c.ShouldBindJSON(&newUser); err != nil{
			c.JSON(http.StatusBadRequest, err)
			return
		}

		// check if user exists
		exists, err := checkExists(newUser.Email)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if exists{
			c.JSON(http.StatusConflict, gin.H{"data": "user already exists"})
			return
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil{
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": result})
	}
}


func checkExists(email *string) (bool, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": *email}, nil)
	if err != nil {
		return false, err
	}

	if count > 0 {
		// User exists
		return true, nil
	}

	// User does not exist
	return false, nil
}




func UpdateSingleUser() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10 * time.Second)
		defer cancel()

		userID := c.Param("uid")
		id, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Get the field-value pair from the request body
		var updatedData map[string]interface{}
		if err := c.ShouldBindJSON(&updatedData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		update := bson.M{"$set": updatedData}
		updatedUser, err := userCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": updatedUser})
	}
}


func DeleteSingleUser() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10 * time.Second)
		defer cancel()

		user_id := c.Param("uid")   // get param from route
		id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

			// filter to check for email existence
		del, err := userCollection.DeleteOne(ctx, bson.M{"_id": id}, nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		if (del.DeletedCount == 0){
			c.JSON(http.StatusOK, gin.H{"message": "no user instance to delete!!"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": del})
	}
}


func DeleteAllUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10 * time.Second)
		defer cancel()

		idx, err := userCollection.DeleteMany(ctx, bson.D{{}}, nil)
		if err != nil {
			log.Fatal(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{"message": idx })
	}
}