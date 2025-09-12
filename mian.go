package main

import (
	"os"

	"example.com/ecomerce/config"
	"example.com/ecomerce/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
    server := gin.Default()
    
   

		server.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173","https://your-gebeya.vercel.app", "https://your-gebeta.onrender.com"}, // front-end URLs
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, // <-- include Authorization
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
}))

    config.ConnectDatabase()
    routes.RegistorRoutes(server)

    // Use PORT from environment (Render sets this automatically)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // fallback for local development
    }

    server.Run(":" + port)
}
