package config

import (
	"fmt"
	"log"

	"example.com/ecomerce/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
func ConnectDatabase() {
 dsn := "host=dpg-d5kgp4re5dus73a9fnn0-a user=ecomerce2_ps8d_user password=Rm9S8QddlHryCeffjEPWE7SbMqBzw4pS dbname=ecomerce2_ps8d port=5432 sslmode=require"

    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("❌ Failed to connect to database: %v", err)
    }

    DB = db
    fmt.Println("✅ Database connection established")

    db.AutoMigrate(
        &model.Product{},
        &model.User{},
        &model.Image{},
        &model.MerchantProfile{},
        &model.Order{},
        &model.OrderItem{},
    )
}
 func ResetDatabase() {
	
	DB.Migrator().DropTable(
		&model.Product{},
		&model.User{},
		&model.Image{},
		&model.MerchantProfile{},
		&model.Order{},
		&model.OrderItem{},
	)

	println("All tables dropped")

	
	DB.AutoMigrate(
		&model.Product{},
		&model.User{},
		&model.Image{},
		&model.MerchantProfile{},
		&model.Order{},
		&model.OrderItem{},
	)

	println("All tables recreated")
}