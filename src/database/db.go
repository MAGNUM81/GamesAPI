package database

import (
	"GamesAPI/src/utils/errorUtils"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"os"
	"strconv"
)

func Setup(initRepos func(*gorm.DB)) (*gorm.DB, error) {

	var server = os.Getenv("DB_HOST")
	var user = os.Getenv("DB_USERNAME")
	var password = os.Getenv("PASSWORD")
	var databaseName = os.Getenv("DATABASE")

	var port, errPort = strconv.Atoi(os.Getenv("DB_PORT"))
	if errPort != nil {
		return nil, errorUtils.NewError(errPort.Error())
	}

	var db, errDb = Connect(server, port, user, password, databaseName)
	if errDb != nil {
		return nil, errorUtils.NewError(errDb.Error())
	}
	initRepos(db)

	return db, nil
}

func Connect(server string, port int, user string, password string, databaseName string) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;",
		server, user, password, port)
	db, err := gorm.Open("mssql", connectionString)
	if err != nil {
		return nil, err
	}
	createDB := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", databaseName))
	if createDB.Error != nil {
		fmt.Println(createDB.Error)
	}
	return db, nil
}
