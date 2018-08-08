package main

import (
	"github.com/egamorim/purchase-api/application"
)

var (
	port       = ":8000"
	dbUserName = "postgres"
	dbPassword = "123"
	dbName     = "postgres"
)

func main() {

	var a = application.App{}
	a.Initialize(
		dbUserName,
		dbPassword,
		dbName)

	a.Run(port)
}
