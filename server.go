package main

import (
	"goproject/actions"
	"goproject/routes"
	"goproject/storage"
	"net/http"
)

func main() {

	var db storage.IStorage = &storage.SqlLiteDb{}
	db.InitializeDatabase()

	var actions actions.IActions = &actions.Actions{Db: db}

	server := routes.SetupRoutes(actions)
	http.ListenAndServe(":8888", server)
}
