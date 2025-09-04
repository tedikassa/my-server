package main

import (
	"example.com/ecomerce/config"
	"example.com/ecomerce/routes"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	server := gin.Default()
	 server.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"}, // your frontend URL
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))
   config.ConnectDatabase()
   routes.RegistorRoutes(server)
	server.Run(":8080")
}
