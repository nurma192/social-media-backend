package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {

}
func Login(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "User Login",
	})
}
func GetUserById(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "User Get",
	})
}
func UpdateUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "User Update",
	})
}
func Current(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "User Current",
	})
}
