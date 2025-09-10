package routes

import (
	"example.com/ecomerce/controller"
	"example.com/ecomerce/middlewere"

	"github.com/gin-gonic/gin"
)


func RegistorRoutes(server *gin.Engine) {
  server.POST("/api/product",middlewere.AuthMiddlewere,controller.AddProduct)
	server.POST("/api/payment/:id",middlewere.AuthMiddlewere,controller.Payment)
   server.GET("/api/product",middlewere.AuthMiddlewere,controller.GetAllProduct)
	server.GET("/api/product/:id",middlewere.AuthMiddlewere,controller.GetProductById)
	server.POST("/api/signup",controller.Signup)
	server.POST("/api/login",controller.Login)
	server.PATCH("/api/update/user/:id",middlewere.AuthMiddlewere,controller.UpdateUser)
	server.POST("/api/webhook/incoming", controller.SantimpayWebhookIncoming)
	server.POST("/api/merchant/delivery/:id",middlewere.AuthMiddlewere,controller.ConfirmDelivery)
	server.GET("/api/orders",middlewere.AuthMiddlewere,controller.GetAllOrder)
	server.GET("/api/merchant/product/:id",middlewere.AuthMiddlewere,controller.GetMerchantProduct)
	server.GET("/api/merchant/order/:id",middlewere.AuthMiddlewere,controller.MerchantOrder)
	server.POST("/api/webhook/payout",controller.SantimpayWebhookPayout)
	server.POST("/api/merchant/payment/:id",middlewere.AuthMiddlewere,controller.AskPayout)
}
