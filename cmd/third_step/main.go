package main

import "github.com/wralith/transaction-poc/choreography"

func main() {
	config := choreography.GetThirdStepAppConfig()
	inbox := map[string]*choreography.CreateTargetAudiencePatternRequest{}
	db := map[string]*choreography.TargetAudiencePattern{}

	app := choreography.NewThirdStepApp(inbox, db, config)

	app.Listen(config.Port)
}
