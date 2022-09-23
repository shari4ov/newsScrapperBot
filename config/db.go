package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DB struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	DBType     string
}

func returnDbENV() DB {
	var db DB
	err := godotenv.Load("local.env")
	fmt.Println(err)
	db.DBUser = os.Getenv("DBUser")
	db.DBPassword = os.Getenv("DBPassword")
	db.DBName = os.Getenv("DBName")
	db.DBHost = os.Getenv("DBHost")
	db.DBPassword = os.Getenv("DBPassword")
	db.DBPort = os.Getenv("DBPort")
	db.DBType = os.Getenv("DBType")
	return db
}

var db DB = returnDbENV()
var (
	DBUser     = db.DBUser
	DBPassword = db.DBPassword
	DBName     = db.DBName
	DBHost     = db.DBHost
	DBPort     = db.DBPort
	DBType     = db.DBType
)

func GetDBType() string {
	return DBType
}
func GetPostgresConnection() string {
	dataBase := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		DBHost,
		DBPort,
		DBUser,
		DBName,
		DBPassword)
	return dataBase
}
