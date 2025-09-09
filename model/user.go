package model

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type RegisterInput struct {
    Name        string `json:"name" binding:"required,min=3"`
    Email       string `json:"email" binding:"required,email"`
    Password    string `json:"password" binding:"required,min=4"`
    Role        string `json:"role" binding:"required,oneof=user merchant admin"`
    Phone       string `json:"phone" binding:"omitempty"`
}

type User struct {
    gorm.Model
    Name     string `json:"name" `
    Email    string `json:"email" gorm:"unique" `
    Password string `json:"password"`
    Role     string `gorm:"default:'user'"` // "user" or "merchant"
    Phone    string `json:"phone" `

   MerchantProfile MerchantProfile `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"merchantProfile"`
}
type MerchantProfile struct {
    gorm.Model
    UserID      uint    // FK → links to User
    Phone    string `json:"phone"`
    Products []Product `gorm:"foreignKey:MerchantProfileID" json:"products,omitempty"`
}

type Login struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=4"`
}
type Claims struct{
	ID int
	Name string
	Role string
	jwt.RegisteredClaims
}
type UpdateUser struct {
    Name     string `json:"name" binding:"omitempty,min=3"`
    Email    string `json:"email" binding:"omitempty,email"`
    Password string `json:"password" binding:"omitempty,min=6"`
    Phone    string `json:"phone" binding:"omitempty"`
		SantimpayID string `json:"santimpay_id" binding:"omitempty"`
    PrivateKey  string `json:"private_key" binding:"omitempty"`
}
type UpdateMerchantProfile struct
{
	SantimpayID string `json:"santimpay_id" binding:"omitempty"`
  PrivateKey  string `json:"private_key" binding:"omitempty"`
}
