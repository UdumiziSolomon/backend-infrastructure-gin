package user 

import(
	"context"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/solonode/golang-jwt-mongo/internal/database"

)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "user")

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		c.JSON(200, gin.H{"success": "message"})
	}
}

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context){
		var newUser User  

		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		// validate the request body
		if err := c.ShouldBindJSON(&newUser); err != nil{
			c.JSON(http.StatusBadRequest, err)
			return
		}

		// check if user exists
		exists, err := checkExists(newUser.Email)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
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

		c.JSON(http.StatusCreated, result)
	}
}

func checkExists(email *string) (bool, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// filter to check for email existence
	filter := bson.M{"email": email}

	// count the number of matching doc
	var user User
	if err := userCollection.FindOne(ctx, filter).Decode(&user); err != nil{
		if err == mongo.ErrNoDocuments {
			// User does not exist
			return false, nil
		}
		return false, err
	}

	return true, nil
}