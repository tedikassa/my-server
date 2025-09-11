package controller

import (
	"errors"
	"net/http"
	"strconv"

	"example.com/ecomerce/config"
	"example.com/ecomerce/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddProduct(context *gin.Context) {
	var product model.Product
	err:=context.ShouldBindJSON(&product)
	if err!=nil{
		context.JSON(http.StatusBadRequest,gin.H{"status":"fail","message":err.Error()})
		return
	}
  if err=config.DB.Create(&product).Error;err!=nil{
   context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","message":"internal server error,please try again"})
	 return
	}
 context.JSON(http.StatusCreated,gin.H{"status":"success","data":product})

}
func GetAllProduct(context *gin.Context)  {
	var products []model.Product
   if err:=config.DB.Model(&model.Product{}).Preload("Images").Find(&products).Error;err!=nil{
  context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","message":"internal server error,please try again"})
	 return
	 }
	  context.JSON(http.StatusOK,gin.H{"status":"success","data":products})
}

func GetProductById (context *gin.Context){
	var product model.Product
	productId,err:=strconv.Atoi(context.Param("id"))
	if err!=nil {
		context.JSON(http.StatusBadRequest,gin.H{"status":"fail","message":"invaid id"})
		return
	}
	result:=config.DB.Model(&model.Product{}).Preload("Images").First(&product,productId)
	if result.Error!=nil{
		if errors.Is(result.Error,gorm.ErrRecordNotFound){
			context.JSON(http.StatusBadRequest,gin.H{"status":"fail","message":"product not found"})
			return
		}
		context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","message":"server error"})
     return
	}
	 context.JSON(http.StatusCreated,gin.H{"status":"success","data":product})
}

func UpdateProduct(context *gin.Context)  {
	var input model.UpdateProduct
	err:=context.ShouldBindJSON(&input)
	if err!=nil{
			context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","message":"bad data"})
     return
	}
	productId,err:=strconv.Atoi(context.Param("id"))
 if err!=nil {
		context.JSON(http.StatusBadRequest,gin.H{"status":"fail","message":"invaid id"})
		return
	}

	  var product model.Product
    if err := config.DB.First(&product,productId ).Error; err != nil {
        context.JSON(http.StatusNotFound, gin.H{"status":"fail","error": "Product not found"})
        return
    }

    
    if err=config.DB.Model(&product).Updates(input).Error;err!=nil{
			context.JSON(http.StatusInternalServerError, gin.H{"status":"fail","error": "server error"})
        return
		}

    context.JSON(http.StatusOK, gin.H{"status": "success", "product": product})
	 
}
func DeleteProduct(c *gin.Context) {
    productId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "invalid id"})
        return
    }

    // optional: check existence first
    var product model.Product
    if err := config.DB.First(&product, productId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"status": "fail", "error": "Product not found"})
        return
    }

    // now delete
    if err := config.DB.Delete(&product).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "error": "server error"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "success", "product": product})
}

func GetProductsByCategory(c *gin.Context) {
    var products []model.Product
    category := c.Param("category")

    if err := config.DB.Model(&model.Product{}).
        Preload("Images").
        Where("category = ?", category).
        Find(&products).Error; err != nil {

        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "fail",
            "error":  "server error",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data":   products,
    })
}
func GetMerchantProduct(context *gin.Context)  {
	var merchant model.MerchantProfile
	merchantID,_:=strconv.Atoi(context.Param("id"))
	if err:=config.DB.Preload("Products").First(&merchant,merchantID).Error;err!=nil{
			context.JSON(http.StatusNotFound, gin.H{"status":"fail","error": "not found"})
        return
	}
	 context.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data":   merchant,
    })
}


