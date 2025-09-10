package model

import "gorm.io/gorm"

type Order struct {
    gorm.Model
		Key string `gorm:"unique"`
		TransactionID string  
    UserID     uint     `json:"userId" binding:"required"`    
    User      User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Status      string       `json:"status" gorm:"default:'pending'" binding:"omitempty,oneof=pending paid delivered cancelled"`
		
		TotalPrice float64      `json:"totalPrice" binding:"required,gt=0"`
		OrderItems []OrderItem   `gorm:"foreignKey:OrderID"`
}
type OrderItem struct {
	  
    gorm.Model
		PayoutID string 
    OrderID    uint 
		Order      Order  `json:"order,omitempty"`    
    ProductID  uint    `json:"productId" binding:"required"`
    Product    Product `json:"product,omitempty"` // preload if needed
    MerchantProfileID uint    `json:"merchantProfileId" binding:"required"`
		MerchantProfile   `json:"merchantProfile,omitempty"`
    Quantity   int     `json:"quantity" binding:"required,gt=0"`
    Price      float64 `json:"price" binding:"required,gt=0"` 
		Name string `json:"name"`
		Email string `json:"email"`
		Address string `json:"address"`
		DeliveredCode string
		Delivered bool `gorm:"defualt:false"`
		MerStatus bool `gorm:"defualt:false"`
}
type OrderItemResponse struct {
    ID               uint    `json:"id"`
    OrderID          uint    `json:"orderId"`
    OrderStatus      string  `json:"orderStatus"`
		ProductName      string   `json:"productName"`
    ProductID        uint    `json:"productId"`
    MerchantProfileID uint   `json:"merchantProfileId"`
    Quantity         int     `json:"quantity"`
    Price            float64 `json:"price"`
    Name             string  `json:"name"`
    Email            string  `json:"email"`
    Address          string  `json:"address"`
    DeliveredCode    string  `json:"deliveredCode"`
    Delivered        bool    `json:"delivered"`
    MerchantStatus   bool  `json:"merchantStatus"`
}

type SantimWebhook struct {
    ID string `json:"id"`
    TxnID   string `json:"txnId"`
    Status  string `json:"Status"`
    Amount  string `json:"amount"`
}

type ConfirmDeliveryRequest struct {
    ItemID        uint   `json:"itemId" binding:"required"`
    DeliveredCode string `json:"deliveredCode" binding:"required"`
}
type AskPayout struct{
	 ItemID        uint   `json:"itemId" binding:"required"`
}
