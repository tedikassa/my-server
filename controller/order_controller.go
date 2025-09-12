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
func UserOrder(context *gin.Context) {
    userID, _ := strconv.Atoi(context.Param("id"))

    var orders []model.Order
    if err := config.DB.Model(&model.Order{}).
        Preload("OrderItems.Product.Images").
        Where("user_id =? AND status=?", userID,"paid").
        Find(&orders).Error; err != nil {
        context.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
        return
    }

    // Flatten to only necessary info
    var orderItems []map[string]interface{}
    for _, order := range orders {
        for _, item := range order.OrderItems {
            orderItems = append(orderItems, map[string]interface{}{
                "id":           item.ID,
                "orderId":      order.ID,
                "name":         item.Name,
                "price":        item.Price,
                "quantity":     item.Quantity,
                "delivered":    item.Delivered,
                "address":      item.Address,
                "DeliveredCode": item.DeliveredCode,
                "image":        item.Product.Images, 
				"productName":item.Product.Name,
            })
        }
    }

    context.JSON(http.StatusOK, gin.H{"status": "success", "data": orderItems,"orders":orders})
}


func GetAllOrder(context *gin.Context)  {
	var orders []model.Order
	if err:=config.DB.Preload("OrderItems").Find(&orders).Error;err!=nil{
		context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","message":err.Error()})
		return
	}
	context.JSON(http.StatusOK,gin.H{"status":"sucess","data":orders})
}