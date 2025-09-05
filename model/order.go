package model

import "gorm.io/gorm"

type Order struct {
    gorm.Model
    UserID     uint         `json:"buyer_id" binding:"required"`
    User      User         `json:"buyer,omitempty"` 
    Status      string       `json:"status" gorm:"default:'pending'" binding:"omitempty,oneof=pending paid shipped delivered cancelled"`
		ProductID  uint    `json:"product_id" binding:"required"`
    Product    Product `json:"product,omitempty"` 
		Quantity   int     `json:"quantity" binding:"required,gt=0"`
		MerchantID uint    `json:"merchant_id" binding:"required"`
		Price      float64 `json:"price" binding:"required,gt=0"` 
		TotalPrice float64      `json:"total_amount" binding:"required,gt=0"`
}
type SantimWebhook struct {
    TxnID   string `json:"txnId"`
    Status  string `json:"Status"`
    Amount  string `json:"amount"`
    
}

// OrderItem represents a single product in an order.
// type OrderItem struct {
//     gorm.Model
//     OrderID    uint    `json:"order_id"` // FK → links to Order
//     ProductID  uint    `json:"product_id" binding:"required"`
//     Product    Product `json:"product,omitempty"` // preload if needed
//     MerchantID uint    `json:"merchant_id" binding:"required"`
//     Quantity   int     `json:"quantity" binding:"required,gt=0"`
//     Price      float64 `json:"price" binding:"required,gt=0"` // snapshot of product price at order time
// }
