package main

import (
	"fmt"
	"net/smtp"
	"log"
	"time"
	"math/rand"
)


func main(){
	fmt.Println("Trying to send mail...")
	sendVerificationCode()
}

func sendVerificationCode(){
	// code
	rand.Seed(time.Now().UnixNano())
	verificationCode := rand.Intn(900000) + 100000 

	//  sender data
	senderMail := "solomonudumizi@gmail.com"
	senderPassword := "nnrfihcqmkqmrohq"
	toAddr := "solonode21@gmail.com"

	// SMTP options
	smtpHost := "smtp.gmail.com"
	smtpPort :=  465

	// message
	subject := "Golang emailing process"
	body := fmt.Sprintf("<h1>First email messaging</h1>,<h2> verification code: </h2>", verificationCode)

	fmt.Println("Tring to send a mail!!!........")

	m := gomail.NewMessage()
	m.SetHeader("From", senderMail)
	m.SetHeader("To", toAddr)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	mail := gomail.NewDialer(smtpHost, smtpPort, senderMail, senderPassword)

	if err := mail.DialAndSend(m); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mail was sent!!")
}

func VerifyUser() gin.HandlerFunc {
	return(c *gin.Context){
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
		body := fmt.Sprintf("<h1> %v </h1>", verificationCode)

		m := gomail.NewMessage()  // new instance for gomail
		m.SetHeader("From", config.MailSender)
		m.SetHeader("To", toAddr)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)

		mail := gomail.NewDialer(config.SmtpHost, config.SmtpPort, config.MailSender, config.MailPassword)
		if err := mail.DialAndSend(m); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
		}

		c.JSON(http.StatusOK, gin.H{"msg": "Mail was sent"})

	}
}

