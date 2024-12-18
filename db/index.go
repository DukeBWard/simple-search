package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBConnect *gorm.DB

func InitDB() {
	dburl := os.Getenv("DATABASE_URL")
	var err error
	DBConnect, err = gorm.Open(postgres.Open(dburl))
	if err != nil {
		fmt.Println("failed to connect to DB")
		panic(err)
	}

	// uuid extension
	err = DBConnect.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		fmt.Println("cannot install uuid")
		panic(err)
	}
	err = DBConnect.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return DBConnect
}
