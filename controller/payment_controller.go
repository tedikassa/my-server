package controller

import (
	"encoding/hex"
	"fmt"
	"time"

	"math/rand"
	"net/http"
	"strconv"

	"example.com/ecomerce/config"
	"example.com/ecomerce/model"
	"example.com/ecomerce/sdk"
	"github.com/gin-gonic/gin"
)

func Payment(context *gin.Context) {
	var order model.Order
	userId:=context.Param("id")
	err:=context.ShouldBindJSON(&order)
	if err!=nil {
		 	context.JSON(http.StatusBadRequest,gin.H{"status":"fail","error":err.Error(),})
	   return
	}
	
 
	// Santim Test
const PRIVATE_KEY_IN_PEM = `
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIF/mI9tSZxKbfEniC+3yfvwIS/D76+p/ky/oDmKAwu5roAoGCCqGSM49
AwEHoUQDQgAEqJl+TIowE6CAhoghgmH+cdzn5+WNax9/REqXJf6b1HdJCRZBCXWT
6coLZ23OyF5x9uVOUXixZeB7J7y9iSWDzw==
-----END EC PRIVATE KEY-----
`

	const GATEWAY_MERCHANT_ID = "9e2dab64-e2bb-4837-9b85-d855dd878d2b"

	const testBed = false

	sdk, err := sdk.NewSantimpaySDK(GATEWAY_MERCHANT_ID, PRIVATE_KEY_IN_PEM, testBed)
	if err != nil {
		context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","error":"server error"})
		return
	}

	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(1000000000)
	strid := strconv.Itoa(int(id))
	uId,_:=strconv.Atoi(userId)
	order.UserID=uint(uId)
   order.TransactionID=strid
	  for i := range order.OrderItems {
        order.OrderItems[i].DeliveredCode = generateCode(4) // 8 hex chars
    }
 if err=config.DB.Create(&order).Error;err!=nil{
		context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","error":"server error"})
	   return
 }

	const phoneNumber = ""//"+251909090909"
	const notifyURL = "https://your-gebeta.onrender.com/api/webhook"
	const successRedirectURL = "https://santimpay.com"
	const failureRedirectURL = "https://santimpay.com"
	const cancelRedirectURL = "https://santimpay.com"

	// Generate a payment URL
	paymentURL, err := sdk.GeneratePaymentURL(strid, order.TotalPrice, "online market", successRedirectURL, failureRedirectURL, notifyURL, phoneNumber, cancelRedirectURL)
	if err != nil {
			context.JSON(http.StatusUnauthorized,gin.H{"status":"fail","error":"please enter valid private key or merchantId"})
			return
	}

	fmt.Println("Payment URL:", paymentURL)
  context.JSON(http.StatusOK,gin.H{"paymenturl":paymentURL,"data":order})

}



func SantimpayWebhook(c *gin.Context) {
    var payload model.SantimWebhook

    if err := c.BindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "invalid JSON"})
        return
    }

    // Convert amount string to float64 if you want
    amount, _ := strconv.ParseFloat(payload.Amount, 64)

    fmt.Printf("Transaction ID: %s, Status: %s, Amount: %.2f\n",
        payload.TxnID, payload.Status, amount)
     payload.Status = "SUCCESS"
    // Example: update your order
    if payload.Status == "SUCCESS" {
        config.DB.Model(&model.Order{}).
            Where("transaction_id = ?", payload.TxnID).
            Updates(map[string]interface{}{
                "status": "paid",
                "total_price": amount,
            })
    }

    c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func generateCode(n int) string {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return "default"
    }
    return hex.EncodeToString(b)
}


/*
git add -A
git commit -m 'update'
git push -u origin main 

*/