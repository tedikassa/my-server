package middlewere

import (
	"net/http"
	"strings"

	"example.com/ecomerce/model"
	"example.com/ecomerce/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddlewere(context *gin.Context) {
	authHeader := context.GetHeader("Authorization");
	if authHeader==""{
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status":"fail","error": "please log in again"})
        return
	}
	var token string
	tokenArray:=strings.SplitN(authHeader," ",2)
	if len(tokenArray)==2&&tokenArray[0]=="Bearer"{
   token=tokenArray[1]
	}else{
		token=authHeader
	}
	claims:=&model.Claims{}
	parsedToken,err:=jwt.ParseWithClaims(token,claims,func(t *jwt.Token) (any, error) {
		return utils.JwtKey,nil
	})
	if err!=nil||!parsedToken.Valid{
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status": "fail",
				"error":  "invalid or expired token",})
	}
  context.Set("id",claims.ID)
	context.Set("name",claims.Name)
	context.Set("role",claims.Role)
	
}

func MerchantMiddleware(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists || role != "merchant" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status": "fail",
			"error":  "access restricted to merchants only",
		})
		return
	}
	c.Next()
}
