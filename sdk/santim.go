package sdk

import (
	// "context"
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"net/http"

	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SantimpaySdk struct {
	MerchantID string
	PrivateKey *ecdsa.PrivateKey
	BaseURL    string
	HTTPClient *http.Client
}

const PRODUCTION_BASE_URL = "https://services.santimpay.com/api/v1/gateway"
const TEST_BASE_URL = "https://testnet.santimpay.com/api/v1/gateway"

func NewSantimpaySDK(MerchantID string, PrivateKey string, testBed bool) (*SantimpaySdk, error) {
	baseURL := PRODUCTION_BASE_URL
	if testBed {
		baseURL = TEST_BASE_URL
	}

	block, _ := pem.Decode([]byte(PrivateKey))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &SantimpaySdk{
		MerchantID: MerchantID,
		PrivateKey: key,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}, nil
}

func sign(payload map[string]interface{}, privateKey *ecdsa.PrivateKey, signingMethod jwt.SigningMethod) (string, error) {
	token := jwt.NewWithClaims(signingMethod, jwt.MapClaims(payload))

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (sdk *SantimpaySdk) generateToken(payload map[string]interface{}, signingMethod jwt.SigningMethod) (string, error) {
	return sign(payload, sdk.PrivateKey, signingMethod)
}

func (sdk *SantimpaySdk) generateSignedTokenForInitialPayment(amount float64, paymentReason string) (string, error) {
	payload := map[string]interface{}{
		"amount":        amount,
		"paymentReason": paymentReason,
		"merchantID":    sdk.MerchantID,
		"generated":     time.Now().Unix(),
	}
	// fmt.Print(sdk.generateToken(payload, jwt.SigningMethodES256))
	return sdk.generateToken(payload, jwt.SigningMethodES256)
}

func (sdk *SantimpaySdk) generateSignedTokenForDirectPayment(amount float64, paymentReason string, paymentMethod string, phoneNumber string) (string, error) {
	payload := map[string]interface{}{
		"amount":        amount,
		"paymentReason": paymentReason,
		"paymentMethod": paymentMethod,
		"phoneNumber":   phoneNumber,
		"merchantID":    sdk.MerchantID,
		"generated":     time.Now().Unix(),
	}

	return sdk.generateToken(payload, jwt.SigningMethodES256)
}

func (sdk *SantimpaySdk) generateSignedTokenForGetTransaction(id string) (string, error) {
	time := time.Now().Unix()
	claims := jwt.MapClaims{
		"id":        id,
		"merId":     sdk.MerchantID,
		"generated": time,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedToken, err := token.SignedString(sdk.PrivateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (sdk *SantimpaySdk) GeneratePaymentURL(id string, amount float64, paymentReason, successRedirectURL, failureRedirectURL, notifyURL, phoneNumber, cancelRedirectURL string) (string, error) {
	token, err := sdk.generateSignedTokenForInitialPayment(amount, paymentReason)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return "", err
	}

	payload := map[string]interface{}{
		"id":                 id,
		"amount":             amount,
		"reason":             paymentReason,
		"merchantID":         sdk.MerchantID,
		"signedToken":        token,
		"successRedirectURL": successRedirectURL,
		"failureRedirectURL": failureRedirectURL,
		"notifyURL":          notifyURL,
		"cancelRedirect":     cancelRedirectURL,
	}

	if phoneNumber != "" {
		payload["phoneNumber"] = phoneNumber
	}

	jsonPayload, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonPayload))

	resp, err := http.Post(sdk.BaseURL+"/initiate-payment", "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return "", err
	}

	// fmt.Printf("response: %+v\n", resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	if resp.StatusCode == http.StatusOK {
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Println("Error decoding response body:", err)
			return "", err
		}
		url, ok := response["url"].(string)
		if !ok {
			fmt.Println("URL not found in the response")
			return "", errors.New("URL not found in the response")
		}
		return url, nil
	} else {
		 fmt.Println("HTTP Request Error:", resp.Status)
		 fmt.Println("Response Body:", string(body))
		 fmt.Println("Response Headers:", resp.Header)
		return "", errors.New("failed to initiate payment")
	}
}

func (sdk *SantimpaySdk) SendToCustomer(id string, amount float64, paymentReason, phoneNumber, paymentMethod ,notifyUrl string) (interface{}, error) {
	token, err := sdk.generateSignedTokenForDirectPayment(amount, paymentReason, paymentMethod, phoneNumber)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"id":                    id,
		"clientReference":       id,
		"amount":                amount,
		"reason":                paymentReason,
		"merchantId":            sdk.MerchantID,
		"signedToken":           token,
		"receiverAccountNumber": phoneNumber,
		"paymentMethod":         paymentMethod,
		"notifyUrl": notifyUrl,
		 
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	
	resp, err := sdk.HTTPClient.Post(sdk.BaseURL+"/payout-transfer", "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		fmt.Println("error:",err.Error());
		return nil, err
	}
body, _ := io.ReadAll(resp.Body)
fmt.Println("status code:", resp.StatusCode)
fmt.Println("response body:", string(body))
if resp.StatusCode == http.StatusOK {
    var response map[string]interface{}
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, err
    }
    return response, nil 
}else {
    return nil, fmt.Errorf("failed to initiate B2C: %s", string(body))
}
}

func (sdk *SantimpaySdk) generateSignedTokenForDirectPaymentOrB2C(amount float64, paymentReason, paymentMethod, phoneNumber string) (string, error) {
	payload := map[string]interface{}{
		"amount":        amount,
		"paymentReason": paymentReason,
		"paymentMethod": paymentMethod,
		"phoneNumber":   phoneNumber,
		"merchantID":    sdk.MerchantID,
		"generated":     time.Now().Unix(),
	}
	return sdk.generateToken(payload, jwt.SigningMethodES256)
}

func (sdk *SantimpaySdk) directPayment(id string, amount float64, paymentReason, notifyURL, phoneNumber, paymentMethod string) (interface{}, error) {
	token, err := sdk.generateSignedTokenForDirectPayment(amount, paymentReason, paymentMethod, phoneNumber)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"id":          id,
		"amount":      amount,
		"reason":      paymentReason,
		"merchantId":  sdk.MerchantID,
		"signedToken": token,
		"phoneNumber": phoneNumber,
		"notifyUrl":   notifyURL,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	resp, err := sdk.HTTPClient.Post(sdk.BaseURL+"/direct-payment", "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Print("Error decoding response")
			return "", err
		}
		url, ok := response["url"].(string)
		if !ok {
			return "", errors.New("URL not found in the response")
		}
		return url, nil
	} else {
		return nil, errors.New("failed to initiate direct payment")
	}
}

func (sdk *SantimpaySdk) checkTransactionStatus(id string) (interface{}, error) {
	token, err := sdk.generateSignedTokenForGetTransaction(id)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signed token: %v", err)
	}

	payload := map[string]interface{}{
		"id":          id,
		"merchantID":  sdk.MerchantID,
		"signedToken": token,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload to JSON: %v", err)
	}
	resp, err := sdk.HTTPClient.Post(sdk.BaseURL+"/fetch-transaction-status", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	//   fmt.Println("Response Status:", resp.Status)
	fmt.Println("Transaction", string(body))
	//   fmt.Println("id", id)

	if resp.StatusCode == http.StatusOK {
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			return "", fmt.Errorf("failed to decode response body: %v", err)
		}
		url, ok := response["url"].(string)
		if !ok {
			return "", errors.New("URL not found in the response")
		}
		return url, nil
	} else {
		return nil, fmt.Errorf("failed to fetch transaction status: %s", string(body))
	}
}
