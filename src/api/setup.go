package api

import (
	"GamesAPI/src/database"
	"GamesAPI/src/entityBinder"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

const server = "database"
const port = 1433
const user = "sa"
const password = "P4tate!!"
const databaseName = "GamesGoDB"

func ConnectDatabase() {
	database.Instance = database.Connect(server, port, user, password, databaseName)
	entityBinder.BindEntities(database.Instance)
}

