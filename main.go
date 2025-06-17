package main

import (
	"book/model"
	"book/routes"
	"book/databases"
)
// i bright
func main(){
	db := databases.Connect()
	model.ModelMigration(db)
	routes.Routes(db)
}




