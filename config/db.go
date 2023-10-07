package config

import (
	"fmt"
	"os"
)

func GetMySqlDSN() string {
	db_name := os.Getenv("DB_NAME")
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_parameter := "charset=utf8mb4&parseTime=True&loc=Local&tls=false"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", db_user, db_password, db_host, db_port, db_name, db_parameter)
	return dsn
}
