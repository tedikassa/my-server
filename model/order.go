package model

import "gorm.io/gorm"

type Order struct {
    gorm.Model
		TransactionID string `gorm:"unique"` 
    UserID     uint     `json:"userId" binding:"required"`    
    User      User         `json:"user,omitempty"` 
    Status      string       `json:"status" gorm:"default:'pending'" binding:"omitempty,oneof=pending paid delivered cancelled"`
		
		TotalPrice float64      `json:"totalPrice" binding:"required,gt=0"`
		OrderItems []OrderItem   `gorm:"foreignKey:OrderID"`
}
type OrderItem struct {
    gorm.Model
    OrderID    uint     
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
}
type SantimWebhook struct {
    TxnID   string `json:"txnId"`
    Status  string `json:"Status"`
    Amount  string `json:"amount"`
}


