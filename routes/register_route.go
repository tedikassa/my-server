package routes

import (
	"example.com/ecomerce/controller"
	"example.com/ecomerce/middlewere"

	"github.com/gin-gonic/gin"
)


func RegistorRoutes(server *gin.Engine) {
  server.POST("/api/product",middlewere.AuthMiddlewere,middlewere.MerchantMiddleware,controller.AddProduct)
	server.POST("/api/product/update/:id",middlewere.AuthMiddlewere,middlewere.MerchantMiddleware,controller.UpdateProduct)
	server.DELETE("/api/product/delete/:id",middlewere.AuthMiddlewere,middlewere.MerchantMiddleware,controller.DeleteProduct)
	server.POST("/api/payment/:id",middlewere.AuthMiddlewere,controller.Payment)
   server.GET("/api/product",middlewere.AuthMiddlewere,controller.GetAllProduct)
	server.GET("/api/product/:id",middlewere.AuthMiddlewere,controller.GetProductById)
	server.POST("/api/signup",controller.Signup)
	server.POST("/api/login",controller.Login)
	server.PATCH("/api/update/user/:id",middlewere.AuthMiddlewere,controller.UpdateUser)
	server.POST("/api/webhook/incoming", controller.SantimpayWebhookIncoming)
	server.POST("/api/merchant/delivery/:id",middlewere.AuthMiddlewere,middlewere.MerchantMiddleware,controller.ConfirmDelivery)
	server.GET("/api/orders",middlewere.AuthMiddlewere,controller.GetAllOrder)
	server.GET("/api/merchant/product/:id",middlewere.AuthMiddlewere,middlewere.MerchantMiddleware,controller.GetMerchantProduct)
	server.GET("/api/merchant/order/:id",middlewere.AuthMiddlewere,middlewere.MerchantMiddleware,controller.MerchantOrder)
	server.POST("/api/webhook/payout",controller.SantimpayWebhookPayout)
	server.POST("/api/merchant/payment/:id",middlewere.AuthMiddlewere,controller.AskPayout)
	server.GET("/api/user/order/:id",controller.UserOrder)
}
