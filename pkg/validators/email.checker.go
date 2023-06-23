package main

import (
	"bufio"
	"log"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

func main(){

	// create a var to read user input on the terminal
	termScanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter the email to verify and press enter...\n")

	// check incase for multiple mails
	for termScanner.Scan(){
		email := termScanner.Text()

		fmt.Println("Checking domain server......")
		exists, err := verifyEmailDomain(email)
		if err != nil{
			fmt.Printf("Error checking email existence: %s\n", err.Error())
			return
		}

		if exists {
			fmt.Printf("Email '%s' exists online.\n", email)
			os.Exit(1)
		} else {
			fmt.Printf("Email '%s' does not exist online.\n", email)
			os.Exit(1)
		}
	}

	// check for error in terminal
	if err := termScanner.Err(); err != nil {
		log.Fatal("Could not read input: %v\n", err)
	}
}

func verifyEmailDomain(domain string) (bool, error) {

	// Extract the domain from the email address
	parts := strings.Split(domain, "@")
	mailDomain := parts[1]

	// Get MX records for the domain
	mxRecords, err := net.LookupMX(mailDomain)
	if err != nil {
		return false, err
	}

	// Get the mail server address
	mailServer := mxRecords[0].Host
	fmt.Printf("✔ ✔ MX record found: %v \n", mailServer)
	go showLoader("Checking mail address for eligibility.")

	// Connect to the mail server
	client, err := smtp.Dial(fmt.Sprintf("%s:25", mailServer))
	if err != nil {
		return false, err
	}
	defer client.Close()

	// Set the sender and recipient
	err = client.Mail("example@example.com")
	if err != nil {
		return false, err
	}

	err = client.Rcpt(domain)
	if err != nil {
		return false, nil
	}

	// If the recipient address exists, the server will return a nil error
	return true, nil
}

func showLoader(text string) {
	loaderChars := []string{"|", "/", "-", "\\"} // Define loader characters

	for {
		for _, char := range loaderChars {
			fmt.Printf("%s\r%s", char,text) // Print the loader character and overwrite previous output
			time.Sleep(100 * time.Millisecond) // Wait for 100 milliseconds
		}
	}
}





// func showProgressBar(total int) {
// 	for i := 0; i <= total; i++ {
// 		// Calculate progress percentage
// 		percentage := (i * 100) / total

// 		// Generate the progress bar
// 		progress := "[" + getProgressBar(percentage) + "] " + fmt.Sprintf("%d%%", percentage)

// 		fmt.Print("\r" + progress) // Print the progress bar and overwrite previous output
// 		time.Sleep(100 * time.Millisecond) // Wait for 100 milliseconds
// 	}

// 	fmt.Println() // Move to the next line after the progress bar is complete
// }

// func getProgressBar(percentage int) string {
// 	const barLength = 30 // Length of the progress bar
// 	completedLength := (barLength * percentage) / 100
// 	remainingLength := barLength - completedLength

// 	completed := repeatStr("=", completedLength)
// 	remaining := repeatStr(" ", remainingLength)

// 	return completed + remaining
// }

// func repeatStr(str string, count int) string {
// 	result := ""
// 	for i := 0; i < count; i++ {
// 		result += str
// 	}
// 	return result
// }