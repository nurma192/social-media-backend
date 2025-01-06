package services

import "fmt"

func Login(username, password string) string {
	return fmt.Sprintf("Hello, %s! This is the Login service.", username)
}

func Register(username, email string) string {
	return fmt.Sprintf("Hello, %s! This is the Register service.", username)
}
