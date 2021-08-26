package email

import (
	"fmt"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendConfirmEmail(email, id string) error {
	from := mail.NewEmail("Bloom Health", "noreply@bloomhealth.app")
	subject := "Please confirm your email address"
	to := mail.NewEmail("Bloom User", email)
	plainTextContent := "Thank you for signing up for bloom. Please confirm your email by visiting https://api.bloomhealth.app/confirm/" + id
	htmlContent := "Thank you for signing up for bloom. Please confirm your email by clicking <a href=\"https://api.bloomhealth.app/confirm/" + id + "\">here.</a>"
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
	from := mail.NewEmail("Bloom Health", "noreply@bloomhealth.app")
	subject := "Password reset"
	to := mail.NewEmail("Bloom User", email)
	plainTextContent := "Please visit https://api.bloomhealth.app/recover?id=" + id + " to reset your password"
	htmlContent := "You may reset your password by clicking <a href=\"https://api.bloomhealth.app/recover?id=" + id + "\">here.</a>"
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
