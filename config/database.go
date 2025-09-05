package config

import (
	"example.com/ecomerce/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
 func ConnectDatabase()  {
	var err error
	dsn:="host=dpg-d2svpvbe5dus73db7f20-a user=ecomerce_3rsi_user password=N4iWNLISUHvesVIbjIa49ogVvOUr263R dbname=ecomerce_3rsi port=5432 sslmode=disable"
	DB,err=gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil {
		panic("server coudnt connect connect")
	}

 println("Database connection establish")
 DB.AutoMigrate(&model.Product{},&model.User{},&model.Image{},&model.MerchantProfile{},&model.Order{},&model.OrderItem{})
 }