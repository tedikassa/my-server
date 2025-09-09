package routes

import (
	"example.com/ecomerce/controller"

	"github.com/gin-gonic/gin"
)


func RegistorRoutes(server *gin.Engine) {
  server.POST("/api/product",controller.AddProduct)
	server.POST("/api/payment/:id",controller.Payment)
   server.GET("/api/product",controller.GetAllProduct)
	server.GET("/api/product/:id",controller.GetProductById)
	server.POST("/api/signup",controller.Signup)
	server.POST("/api/login",controller.Login)
	server.PATCH("/api/update/user/:id",controller.UpdateUser)
	server.POST("/api/webhook/incoming", controller.SantimpayWebhookIncoming)
	server.POST("/api/merchant/delivery/:id",controller.ConfirmDelivery)
	server.GET("/api/orders",controller.GetAllOrder)
	server.GET("/api/merchant/product/:id",controller.GetMerchantProduct)
	server.GET("/api/merchant/order/:id",controller.MerchantOrder)
	server.POST("/api/webhook/payout",controller.SantimpayWebhookPayout)
	server.POST("/api/merchant/payment/:id",controller.AskPayout)
}
