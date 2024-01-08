package main

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendEmail(emailType string, user User) error {
	from := mail.NewEmail("Auth", os.Getenv("SENDGRID_FROM_EMAIL"))
	to := mail.NewEmail("Example User", user.Email)

	var subject string
	var plainTextContent string
	var htmlContent string

	switch emailType {
	case "confirmationEmail":
		subject = "Confirm Email"
		plainTextContent = "Please confirm your email."
		htmlContent = "<strong>here is your confirmation link</strong>"
	case "forgotPassword":
		subject = "Forgot Password"
		plainTextContent = "You've requested a password reset."
		htmlContent = "<strong>here is your password reset link</strong>"
	default:
		return InvalidMessage("Invalid email type")
	}

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	return nil
}
