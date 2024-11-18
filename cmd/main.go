package main

import (
	"github.com/ItsGabeh/DW-proyecto-tienda/internal/db"
	"github.com/ItsGabeh/DW-proyecto-tienda/internal/routes"
)

func main() {
	router := routes.SetupRouter()
	db.ConnectMongo()
	router.Run(":8080")
}
