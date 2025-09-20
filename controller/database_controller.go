package controller

import (
	"net/http"

	"example.com/ecomerce/config"
	"github.com/gin-gonic/gin"
)

func ResetDatabaseHandler(c *gin.Context) {

	config.ResetDatabase()
	c.JSON(http.StatusOK, gin.H{"message": "Database reset successful"})
}