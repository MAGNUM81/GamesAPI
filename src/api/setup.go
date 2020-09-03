package api

import (
	"GamesAPI/src/database"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

const server = "database"
const port = 1433
const user = "sa"
const password = "P4tate!!"
const databaseName = "GamesGoDB"

func ConnectDatabase() *gorm.DB {
	var db = database.Connect(server, port, user, password, databaseName)
	database.BindEntities(db)
	return db
}

func CreateRepositories() {

}

