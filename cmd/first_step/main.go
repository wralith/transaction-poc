package main

import "github.com/wralith/transaction-poc/choreography"

func main() {
	config := choreography.GetFirstStepAppConfig()
	inbox := map[string]*choreography.CreateProductRequest{}
	db := map[string]*choreography.Product{}

	app := choreography.NewFirstStepApp(inbox, db)

	app.Listen(config.Port)
}
