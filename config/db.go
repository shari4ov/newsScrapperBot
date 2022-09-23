package config

import "fmt"

const (
	DBUser     = "kenan"
	DBPassword = "2001"
	DBName     = "goecho"
	DBHost     = "10.128.0.3"
	DBPort     = "5432"
	DBType     = "postgres"
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
