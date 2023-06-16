package infrastructure

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SqlHandler struct {
	Conn *gorm.DB
}

func Init() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading dotenv: %s", err.Error())
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_USER := os.Getenv("DB_USER")

	DsName := DB_USER + ":" + DB_PASSWORD + "@tcp(mariadb:3306)/" + DB_HOST + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(DsName), &gorm.Config{})
	if err != nil && err.Error() != "" {
		log.Fatal(err.Error())
	}
	fmt.Println("init data base ok")
	sqlHandler := new(SqlHandler)
	sqlHandler.Conn = db
	return db
}
