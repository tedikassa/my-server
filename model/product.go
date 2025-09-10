package model

import "gorm.io/gorm"

type Product struct {
    gorm.Model       
    Name        string   `json:"name" binding:"required,min=3"`   
    Description string   `json:"description" binding:"required,min=10"` 
    Price       float64  `json:"price" binding:"required,gt=0"`      
    Stock       int      `json:"stock" binding:"required,gte=0"`      
    Category    string   `json:"category" binding:"required,min=3"`  
    Images      []Image  `json:"images" gorm:"foreignKey:ProductID"`
    MerchantProfileID uint         `json:"merchant_profile_id" binding:"required"` 
    MerchantProfile   MerchantProfile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"merchant-profile,omitempty"`     
}

type Image struct {
    gorm.Model
    ProductID uint     
    ImageURL  string `json:"url"`  
}
type UpdateProduct struct {
    Name        string   `json:"name" binding:"omitempty,min=3"`       
    Description string   `json:"description" binding:"omitempty,min=10"`
    Price       float64  `json:"price" binding:"omitempty,gt=0"`       
    Stock       int      `json:"stock" binding:"omitempty,gte=0"`           
    Category    string   `json:"category" binding:"omitempty,min=3"`        
}

