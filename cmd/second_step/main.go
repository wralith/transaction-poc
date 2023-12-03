package main

import "github.com/wralith/transaction-poc/choreography"

func main() {
	config := choreography.GetSecondStepAppConfig()
	inbox := map[string]*choreography.CreateSaleOrderRequest{}
	db := map[string]*choreography.SaleOrder{}

	app := choreography.NewSecondStepApp(inbox, db, config)

	app.Listen(config.Port)
}
