package services

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

func (s *AppService) sendMessage(email, header, message string) {
	m := gomail.NewMessage()

	m.SetHeader("From", "nurma192k@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", header)
	m.SetBody("text/plain", message)

	d := gomail.NewDialer("smtp.gmail.com", 587, "nurma192k@gmail.com", "rxnp orjf hucd zftg")

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send email:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
