package user 

import(
	"context"
	"time"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)



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

func HashPassword(password string) (hashedPassword string, err error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		err = fmt.Errorf("could not hash password: %w", err)
		return
	}
	hashedPassword = string(hashedBytes)
	return
}

func VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("password verification failed: %w", err)
	}
	return nil
}
