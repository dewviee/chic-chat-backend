package database

import (
	"chicchat/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(dsn string) *gorm.DB {
	log.Println("Connecting to database")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	log.Println("Connect to database successfully")
	return db
}

func MigratingDatabase(db *gorm.DB) {
	db.AutoMigrate(
		models.Room{},
		models.User{},
	)
}
