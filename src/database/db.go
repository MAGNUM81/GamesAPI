package database

import (
	"GamesAPI/src/errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

var Instance *gorm.DB

func Connect(server string, port int, user string, password string, databaseName string) *gorm.DB {
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;",
		server, user, password, port)
	db, err := gorm.Open("mssql", connectionString)
	if err != nil {
		errors.Panic(err)
	}
	createDB := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", databaseName))
	if createDB.Error != nil {
		fmt.Println(createDB.Error)
	}
	return db
}
