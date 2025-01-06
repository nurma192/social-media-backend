package auth_service

import "fmt"

func Login(username, password string) string {
	return fmt.Sprintf("Hello, %s! This is the Login service.", username)
}
