package email

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendConfirmEmail(email, id string) error {
	from := mail.NewEmail("Bloom Health", os.Getenv("EMAIL_SENDER"))
	subject := "Please confirm your email address"
	to := mail.NewEmail("Bloom User", email)
	plainTextContent := "Thank you for signing up for bloom. Please confirm your email by visiting " + os.Getenv("BACKEND_URL") + "/confirm/" + id
	htmlContent := "Thank you for signing up for bloom. Please confirm your email by clicking <a href=\"" + os.Getenv("BACKEND_URL") + "/confirm/" + id + "\">here.</a>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return nil
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return err
	}
}

func SendRecoveryEmail(email, id string) error {
	from := mail.NewEmail("Bloom Health", os.Getenv("EMAIL_SENDER"))
	subject := "Password reset"
	to := mail.NewEmail("Bloom User", email)
	plainTextContent := "Please visit " + os.Getenv("FRONTEND_URL") + "/recover/setPassword?id=" + id + " to reset your password"
	htmlContent := "You may reset your password by clicking <a href=\"" + os.Getenv("FRONTEND_URL") + "/recover/setPassword?id=" + id + "\">here.</a>. This code will expire in 20 minutes."
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return nil
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return err
	}
}
