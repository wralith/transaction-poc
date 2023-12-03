package choreography

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type (
	CreateTargetAudiencePattern map[string]*CreateTargetAudiencePatternRequest
	TargetAuidencePatternDB     map[string]*TargetAudiencePattern
)

// TargetAudiencePatternService
func NewThirdStepApp(inbox CreateTargetAudiencePattern, db TargetAuidencePatternDB, config ThirdStepAppConfig) *fiber.App {
	app := fiber.New()

	app.Post("/targets", CreateTargetAudiencePatternHandler(inbox))
	app.Post("/targets/rollback/:correlationId", RollbackTargetAudiencePatternHandler(config, inbox))
	app.Post("/targets/commit/:correlationId", CommitTargetAudiencePatternHandler(config, inbox, db))

	return app
}

type TargetAudiencePattern struct {
	Name    string
	Pattern string
}

type CreateTargetAudiencePatternRequest struct {
	CorrelationId string `json:"correlationId"`
	Name          string `json:"name"`
	Pattern       string `json:"pattern"`
}

func CreateTargetAudiencePatternHandler(inbox CreateTargetAudiencePattern) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateTargetAudiencePatternRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		inbox[req.CorrelationId] = &req

		return c.SendStatus(fiber.StatusAccepted)
	}
}

func CommitTargetAudiencePatternHandler(config ThirdStepAppConfig, inbox CreateTargetAudiencePattern, db TargetAuidencePatternDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var coreelationId = c.Params("correlationId")
		var awaitingRequest = inbox[coreelationId]

		db[coreelationId] = &TargetAudiencePattern{
			Name:    awaitingRequest.Name,
			Pattern: awaitingRequest.Pattern,
		}

		// Not production ready ofc, just for demo, this call won't ensure that the first step is done or second step will be done after sending this request
		resCode, _, err := fasthttp.Post(nil, config.SecondStepAppURL+"/orders/commit/"+coreelationId, nil)
		if err != nil && resCode != fiber.StatusCreated {
			return err
		}

		delete(inbox, coreelationId)

		return c.SendStatus(fiber.StatusCreated)
	}
}

func RollbackTargetAudiencePatternHandler(config ThirdStepAppConfig, inbox CreateTargetAudiencePattern) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var coreelationId = c.Params("correlationId")

		resCode, _, err := fasthttp.Post(nil, config.SecondStepAppURL+"/orders/rollback/"+coreelationId, nil)
		if err != nil && resCode != fiber.StatusNoContent {
			return err
		}

		delete(inbox, coreelationId)

		return c.SendStatus(fiber.StatusNoContent)
	}
}
