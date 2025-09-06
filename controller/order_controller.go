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
	context.JSON(http.StatusOK,gin.H{"status":"sucess","data":items})
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