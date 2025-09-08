package controller

import (
	"net/http"
	"strconv"

	"example.com/ecomerce/config"
	"example.com/ecomerce/model"
	"github.com/gin-gonic/gin"
)

func MerchantOrder(context *gin.Context) {
	merchantID,_:=strconv.Atoi(context.Param("id"))
	var items []model.OrderItem
	if err:=config.DB.Preload("Order").Where("merchant_profile_id=?",merchantID).Find(&items).Error;err!=nil{
   	context.JSON(http.StatusNotFound,gin.H{"status":"fail","message":err.Error()})
		return
	}
	var resp []model.OrderItemResponse
	for _, i := range items {
    resp = append(resp, model.OrderItemResponse{
        ID: i.ID,
        OrderID: i.OrderID,
        OrderStatus: i.Order.Status, // only one field you need
        ProductID: i.ProductID,
				ProductName:i.Product.Name,
        MerchantProfileID: i.MerchantProfileID,
        Quantity: i.Quantity,
        Price: i.Price,
        Name: i.Name,
        Email: i.Email,
        Address: i.Address,
        DeliveredCode: i.DeliveredCode,
        Delivered: i.Delivered,
        MerchantStatus: i.MerStatus,
    })
}
	context.JSON(http.StatusOK,gin.H{"status":"sucess","data":resp})
}
func UserOrder(context *gin.Context)  {
		userID,_:=strconv.Atoi(context.Param("id"))
			var orders []model.Order
			if err:=config.DB.Model(&model.Order{}).Preload("OrderItems").Where("user_id=?",userID).Find(&orders).Error;err!=nil{
	    context.JSON(http.StatusNotFound,gin.H{"status":"fail","message":err.Error()})
		return
			}		
context.JSON(http.StatusOK,gin.H{"status":"sucess","data":orders})
}

func GetAllOrder(context *gin.Context)  {
	var orders []model.Order
	if err:=config.DB.Preload("OrderItems").Find(&orders).Error;err!=nil{
		context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","message":err.Error()})
		return
	}
	context.JSON(http.StatusOK,gin.H{"status":"sucess","data":orders})
}