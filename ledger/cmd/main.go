package main

import (
	service "ledger/cmd/service"

	"jcgurango.com/ledger/dbmodel"
)

func main() {
	dbmodel.SetupDB()
	service.Run()
}
