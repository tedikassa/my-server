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
const PRIVATE_KEY_IN_PEM = `
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIF/mI9tSZxKbfEniC+3yfvwIS/D76+p/ky/oDmKAwu5roAoGCCqGSM49
AwEHoUQDQgAEqJl+TIowE6CAhoghgmH+cdzn5+WNax9/REqXJf6b1HdJCRZBCXWT
6coLZ23OyF5x9uVOUXixZeB7J7y9iSWDzw==
-----END EC PRIVATE KEY-----
`

	const GATEWAY_MERCHANT_ID = "9e2dab64-e2bb-4837-9b85-d855dd878d2b"

	const testBed = false

func Payment(context *gin.Context) {
	var order model.Order
	err:=context.ShouldBindJSON(&order)
	if err!=nil {
		 	context.JSON(http.StatusBadRequest,gin.H{"status":"fail","error":err.Error(),})
	   return
	}
	

	sdk, err := sdk.NewSantimpaySDK(GATEWAY_MERCHANT_ID, PRIVATE_KEY_IN_PEM, testBed)
	if err != nil {
		context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","error":"server error"})
		return
	}
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(1000000000)
 strid := strconv.Itoa(int(id))
order.Key=strid
	
	  for i := range order.OrderItems {
        order.OrderItems[i].DeliveredCode = generateCode(4) // 8 hex chars
    }
 if err=config.DB.Create(&order).Error;err!=nil{
		context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","error":"server error"})
	   return
 }


	const phoneNumber = ""//"+251909090909"
	const notifyURL = "https://your-gebeta.onrender.com/api/webhook/incoming"
	const successRedirectURL = "https://santimpay.com"
	const failureRedirectURL = "https://santimpay.com"
	const cancelRedirectURL = "https://santimpay.com"

	// Generate a payment URL
	paymentURL, err := sdk.GeneratePaymentURL(strid, order.TotalPrice, "online market", successRedirectURL, failureRedirectURL, notifyURL, phoneNumber, cancelRedirectURL)
	if err != nil {
			context.JSON(http.StatusUnauthorized,gin.H{"status":"fail","error":"please enter valid private key or merchantId"})
			return
	}
  fmt.Println("key:",strid)
	fmt.Println("Payment URL:", paymentURL)
  context.JSON(http.StatusOK,gin.H{"paymenturl":paymentURL,"data":order})

}



func SantimpayWebhookIncoming(c *gin.Context) {
    var payload map[string]interface{} // generic map to see all fields
    if err := c.BindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "fail",
            "error":  "invalid JSON",
        })
        return
    }

    // Debug: log entire payload
    fmt.Printf("SantimPay webhook: %+v\n", payload)
    // Extract fields safely
    txnID, _ := payload["txnId"].(string)
    status, _ := payload["Status"].(string)
    amountStr, _ := payload["amount"].(string)
    key, _ := payload["thirdPartyId"].(string)

    amount, _ := strconv.ParseFloat(amountStr, 64)

    fmt.Printf("TxnID: %s, Status: %s, Amount: %.2f, key: %s\n",
        txnID, status, amount, key)
      status="SUCCESS"
    // Use your own clientReference to update the order
    if status == "SUCCESS" {
        result := config.DB.Model(&model.Order{}).
            Where("key = ?", key).
            Updates(map[string]interface{}{
                "status":      "paid",
                "total_price": amount,
								"transaction_id":txnID,
            })

        fmt.Println("Rows updated:", result.RowsAffected)
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
func ConfirmDelivery(context *gin.Context) {
    
    merchantID, _ := strconv.Atoi(context.Param("id"))

    var req model.ConfirmDeliveryRequest
    if err := context.ShouldBindJSON(&req); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "bad data"})
        return
    }

    // Find the order item
    var item model.OrderItem
    if err := config.DB.First(&item, req.ItemID).Error; err != nil {
        context.JSON(http.StatusNotFound, gin.H{"status": "fail", "error": "order item not found"})
        return
    }
		if item.MerStatus{
				context.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "error": "you have paid"})
        return
		}

    // Check merchant ID
    if item.MerchantProfileID != uint(merchantID) {
        context.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "error": "unauthorized"})
        return
    }
    // Check delivery code
    if item.DeliveredCode != req.DeliveredCode {
        context.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "error": "wrong delivery code"})
        return
    }

    // Update delivered status
    if err := config.DB.Model(&item).Update("delivered", true).Error; err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "error": "server error"})
        return
    }
		sdk, err := sdk.NewSantimpaySDK(GATEWAY_MERCHANT_ID, PRIVATE_KEY_IN_PEM, testBed)
	if err != nil {
		context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","error":"server error"})
		return
	}
	var merchant model.MerchantProfile
if err := config.DB.First(&merchant, item.MerchantProfileID).Error; err != nil {
    context.JSON(http.StatusNotFound, gin.H{"status": "fail", "error": "merchant not found"})
    return
}
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(1000000000)
 strid := strconv.Itoa(int(id))
// phoneNumber:="+251938646985";
 fmt.Println("clientRefernce:",strid)
  notifyURL:= "https://your-gebeta.onrender.com/api/webhook/payout"
	resp,err:=sdk.SendToCustomer(strid,1,"for delivered order","+251906626496", "Telebirr",notifyURL)
	
	if err != nil {
    context.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "error": err.Error()})
    return
		} 
		result,_ := resp.(map[string]interface{})
   payoutID,_:=result["txnId"].(string)
    if err := config.DB.Model(&item).Update("payout_id",payoutID ).Error; err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "error": "server error"})
        return
    }
    // Respond success
    context.JSON(http.StatusOK, gin.H{
        "status": "success",
        "message": "delivery confirmed",
        "payout": resp,
    })
}
   
func AskPayout(context *gin.Context)  {
	merchantID, _ := strconv.Atoi(context.Param("id"))
	var input model.AskPayout
	 if err := context.ShouldBindJSON(&input); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "bad data"})
        return
    }
    // Find the order item
    var item model.OrderItem
    if err := config.DB.First(&item, input.ItemID).Error; err != nil {
        context.JSON(http.StatusNotFound, gin.H{"status": "fail", "error": "order item not found"})
        return
    }

    // Check merchant ID
    if item.MerchantProfileID != uint(merchantID) {
        context.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "error": "unauthorized"})
        return
    }
		if !item.Delivered{
			context.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "error": "order is not delivered"})
        return
		}
		if item.MerStatus{
				context.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "error": "you have paid"})
        return
		}
		sdk, err := sdk.NewSantimpaySDK(GATEWAY_MERCHANT_ID, PRIVATE_KEY_IN_PEM, testBed)
	if err != nil {
		context.JSON(http.StatusInternalServerError,gin.H{"status":"fail","error":"server error"})
		return
	}
	var merchant model.MerchantProfile
if err := config.DB.First(&merchant, item.MerchantProfileID).Error; err != nil {
    context.JSON(http.StatusNotFound, gin.H{"status": "fail", "error": "merchant not found"})
    return
}
		rand.Seed(time.Now().UnixNano())
	id := rand.Intn(1000000000)
 strid := strconv.Itoa(int(id))
// phoneNumber:="+251938646985";
 fmt.Println("clientRefernce:",strid)
  notifyURL:= "https://your-gebeta.onrender.com/api/webhook/payout"
	resp,err:=sdk.SendToCustomer(strid,1,"for delivered order","+251906626496", "Telebirr",notifyURL)
	
	if err != nil {
    context.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "error": err.Error()})
    return
		} 
		result,_ := resp.(map[string]interface{})
   payoutID,_:=result["txnId"].(string)
    if err := config.DB.Model(&item).Update("payout_id",payoutID ).Error; err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "error": "server error"})
        return
    }
    // Respond success
    context.JSON(http.StatusOK, gin.H{
        "status": "success",
        "message": "delivery confirmed",
        "payout": resp,
    })
	
}



func SantimpayWebhookPayout(c *gin.Context)  {
	 var payload map[string]interface{} // generic map to see all fields
    if err := c.BindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "fail",
            "error":  "invalid JSON",
        })
        return
    }

    // Debug: log entire payload
    fmt.Printf("SantimPay webhook: %+v\n", payload)
    // Extract fields safely
    txnID, _ := payload["txnId"].(string)
    status, _ := payload["Status"].(string)
    amountStr, _ := payload["amount"].(string)
    key, _ := payload["thirdPartyId"].(string)

    amount, _ := strconv.ParseFloat(amountStr, 64)

    fmt.Printf("TxnID: %s, Status: %s, Amount: %.2f, key: %s\n",
        txnID, status, amount, key)
    if status == "COMPLETED" {
        result := config.DB.Model(&model.OrderItem{}).
            Where("payout_id = ?", txnID).
            Update("mer_status",true)

        fmt.Println("Rows updated:", result.RowsAffected)
    }

    c.JSON(http.StatusOK, gin.H{"status": "success"})
}

/*
git add -A
git commit -m 'update'
git push -u origin main 

*/