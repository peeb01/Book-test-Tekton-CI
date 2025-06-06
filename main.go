package main

import (
	"book/model"
	"book/routes"
	"book/databases"
)

func main(){
	db := databases.Connect()
	model.ModelMigration(db)
	routes.Routes(db)
}




