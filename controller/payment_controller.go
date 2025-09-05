package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"math/rand"
	"net/http"
	"strconv"
	"time"

	"example.com/ecomerce/sdk"
	"github.com/gin-gonic/gin"
)

func Payment(context *gin.Context) {
//	var order model.Order
	//userId:=context.Param("id")
	// err:=context.ShouldBindJSON(&order)
	// if err!=nil {
	// 	 	context.JSON(http.StatusBadRequest,gin.H{"status":"fail","error":"bad data"})
	//    return
	// }
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
	strid := strconv.Itoa(id)

	const phoneNumber = ""//"+251909090909"
	const notifyURL = "https://your-gebeta.onrender.com/api/webhook"
	const successRedirectURL = "https://santimpay.com"
	const failureRedirectURL = "https://santimpay.com"
	const cancelRedirectURL = "https://santimpay.com"

	// Generate a payment URL
	paymentURL, err := sdk.GeneratePaymentURL(strid, 10, "online market", successRedirectURL, failureRedirectURL, notifyURL, phoneNumber, cancelRedirectURL)
	if err != nil {
			context.JSON(http.StatusUnauthorized,gin.H{"status":"fail","error":"please enter valid private key or merchantId"})
			return
	}

	fmt.Println("Payment URL:", paymentURL)
  context.JSON(http.StatusOK,paymentURL)

}

func SantimpayWebhook(c *gin.Context){

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "cannot read request body"})
		return
	}

	// Log the raw body (optional, for debugging)
	fmt.Println("Webhook payload:", string(body))

	// Parse JSON payload
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "error": "invalid JSON"})
		return
	}

	transactionID, _ := payload["id"].(string)
	status, _ := payload["status"].(string)
	amount, _ := payload["amount"].(float64)

	fmt.Printf("Transaction ID: %s, Status: %s, Amount: %.2f\n", transactionID, status, amount)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

/*
git add -A
git commit -m 'update'
git push -u origin main 

*/