package server

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=B1ntang12!! dbname=licensi_app port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	fmt.Println("Database Connected!")

	// Auto migrate tabel users
	DB.AutoMigrate(&User{})
}
