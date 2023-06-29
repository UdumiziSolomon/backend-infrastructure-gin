package main

import (
	"fmt"
	"math/rand"
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

// type Rectangle struct {   // rectangle struct
// 	width, height  float64
// }

// func area(r Rectangle) float64 {  
// 	return r.width * r.height
// }

// // Method
// func (r Rectangle) Area() float64 {
// 	return r.width * r.height
// }

func main() {

	// instance one
	// rectOne := Rectangle{1,2}
	// rectTwo := Rectangle{30,40}
	// fmt.Printf("%v\n%v\n", area(rectOne), area(rectTwo))

	// //  instance
	// methOne := Rectangle{10,20}
	// fmt.Println(methOne.Area())
	SendVerificationCode()
}

func SendVerificationCode() {
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