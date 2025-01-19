package mail

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendMessage(email, header, message string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "nurma192k@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", header)
	m.SetBody("text/plain", message)

	// todo delete google password and write it to env
	d := gomail.NewDialer("smtp.gmail.com", 587, "nurma192k@gmail.com", "rxnp orjf hucd zftg")

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}
