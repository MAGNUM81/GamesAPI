package database

import (
	"GamesAPI/src/domain"
	"GamesAPI/src/utils/errorUtils"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"os"
	"strconv"
)

func Setup() *gorm.DB {

	var server = os.Getenv("DB_HOST")
	var user = os.Getenv("DB_USERNAME")
	var password = os.Getenv("PASSWORD")
	var databaseName = os.Getenv("DATABASE")

	var port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		errorUtils.Fatal(errorUtils.NewError("Couldn't read db port from .env file."))
	}

	var db = Connect(server, port, user, password, databaseName)
	domain.GameRepo.Initialize(db)
	return db
}

func Connect(server string, port int, user string, password string, databaseName string) *gorm.DB {
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;",
		server, user, password, port)
	db, err := gorm.Open("mssql", connectionString)
	if err != nil {
		errorUtils.Panic(err)
	}
	createDB := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", databaseName))
	if createDB.Error != nil {
		fmt.Println(createDB.Error)
	}
	return db
}
