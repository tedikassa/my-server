package controller

import (
	"net/http"
	"strconv"

	"example.com/ecomerce/config"
	"example.com/ecomerce/model"
	"example.com/ecomerce/utils"
	"github.com/gin-gonic/gin"
)

func Signup(context *gin.Context) {
    var input model.RegisterInput
    if err := context.ShouldBindJSON(&input); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
		var merchant model.MerchantProfile
    hashedPassword, _ := utils.HashPassword(input.Password)

    user := model.User{
        Name:     input.Name,
        Email:    input.Email,
        Password: hashedPassword,
        Role:     input.Role,
        Phone:    input.Phone,
    }

    if err := config.DB.Create(&user).Error; err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
        return
    }
		if input.Role=="merchant"{
			 merchant = model.MerchantProfile{
				UserID: user.ID,
						Phone: user.Phone,
        }
        config.DB.Create(&merchant)
		}    
    context.JSON(http.StatusOK, gin.H{"status": "success", "user": user,"merchant":merchant})
}


func Login(context *gin.Context)  {
	var input model.Login
	err:=context.ShouldBindJSON(&input)
	if err!=nil {
			context.JSON(http.StatusBadRequest, gin.H{"status":"fail","error": "bad data"})
        return
	}
	var user model.User
	if err=config.DB.Preload("MerchantProfile").Where("email=?",input.Email).First(&user).Error;err!=nil{
	  	context.JSON(http.StatusUnauthorized, gin.H{"status":"fail","error": "invalid enail or password"})
        return
	}
	isValid:=utils.CheckHashing(user.Password,input.Password)
	if !isValid{
		context.JSON(http.StatusUnauthorized, gin.H{"status":"fail","error": "invalid enail or password"})
      return
	}
	token,err:=utils.GenerateToken(int(user.ID),user.Name,user.Role)
	if err!=nil{
		context.JSON(http.StatusInternalServerError, gin.H{"status":"fail","error": "could not generate token"})
      return
	}
	context.JSON(http.StatusOK,gin.H{"status":"success","data":user,"token":token})	

}

func UpdateUser(context *gin.Context)  {
	 ID,err:=strconv.Atoi(context.Param("id"))
	 if err!=nil{
		context.JSON(http.StatusBadRequest, gin.H{"status":"fail","error": "invalid id"})
      return
	}
	var input model.UpdateUser
	err=context.ShouldBindJSON(&input)
	if err!=nil{
		context.JSON(http.StatusBadRequest, gin.H{"status":"fail","error": "bad data"})
      return
	}
	var user model.User
	
	if err=config.DB.Preload("MerchantProfile").First(&user,ID).Error;err!=nil{
    context.JSON(http.StatusNotFound, gin.H{"status":"fail","error": "not found"})
      return 
	}
	    if input.Password!=""{
		hashPassword,err:=utils.HashPassword(input.Password)
		if err!=nil{
			context.JSON(http.StatusInternalServerError, gin.H{"status":"fail","error": "could not hash"})
      return 
		}
		input.Password=hashPassword
	 }
	 var merchantProfile model.MerchantProfile
  var updateMerchantProfile model.UpdateMerchantProfile
	 if user.Role=="merchant" {
       merchantProfile=user.MerchantProfile
			 updateMerchantProfile.PrivateKey=input.PrivateKey
			 updateMerchantProfile.SantimpayID=input.SantimpayID
	 }

	 if err=config.DB.Model(&user).Updates(input).Error;err!=nil{
     context.JSON(http.StatusInternalServerError, gin.H{"status":"fail","error": "server error"})
      return 
	 }

	 if user.Role=="merchant"{
		if err=config.DB.Model(&merchantProfile).Updates(updateMerchantProfile).Error;err!=nil{
     context.JSON(http.StatusInternalServerError, gin.H{"status":"fail","error": "server error"})
      return 
	 }
	 }
	
	
	context.JSON(http.StatusOK,gin.H{"status":"success","data":user})
}
func GetUserById(context *gin.Context)  {
  userId,err:=strconv.Atoi(context.Param("id"))
	if err!=nil{
     context.JSON(http.StatusBadRequest, gin.H{"status":"fail","error": "invalid id"})
      return 
	 }
	var user model.User
	if err:=config.DB.First(&user,userId).Error;err!=nil{
		 context.JSON(http.StatusNotFound, gin.H{"status":"fail","error": "not found"})
      return 
	}
  context.JSON(http.StatusOK,gin.H{"status":"success","data":user})
}