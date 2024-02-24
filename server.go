package main

import (
	"goproject/routes"
	"goproject/storage"
	"net/http"
)

func main() {

	var db storage.StorageInterface = &storage.SqlLiteDb{}
	db.InitializeDatabase()

	server := routes.SetupRoutes(db)
	http.ListenAndServe(":8888", server)
}
