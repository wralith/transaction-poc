package choreography

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type (
	CreateSaleOrderInbox map[string]*CreateSaleOrderRequest
	SaleOrderDB          map[string]*SaleOrder
)

// SaleOrderService
func NewSecondStepApp(inbox CreateSaleOrderInbox, db SaleOrderDB, config SecondStepAppConfig) *fiber.App {
	app := fiber.New()

	app.Post("/orders", CreateSaleOrderHandler(inbox))
	app.Post("/orders/rollback/:correlationId", RollbackSaleOrderHandler(config, inbox))
	app.Post("/orders/commit/:correlationId", CommitSaleOrderHandler(config, inbox, db))

	return app
}

type SaleOrder struct {
	Name  string
	Price int
}

type CreateSaleOrderRequest struct {
	CorrelationId string `json:"correlationId"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
}

func CreateSaleOrderHandler(inbox CreateSaleOrderInbox) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateSaleOrderRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		inbox[req.CorrelationId] = &req

		return c.SendStatus(fiber.StatusAccepted)
	}
}

func CommitSaleOrderHandler(config SecondStepAppConfig, inbox CreateSaleOrderInbox, db SaleOrderDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var coreelationId = c.Params("correlationId")
		var awaitingRequest = inbox[coreelationId]

		db[coreelationId] = &SaleOrder{
			Name:  awaitingRequest.Name,
			Price: awaitingRequest.Price,
		}

		// Not production ready ofc, just for demo, this call won't ensure that the first step is done or second step will be done after sending this request
		resCode, _, err := fasthttp.Post(nil, config.FirstStepAppURL+"/products/commit/"+coreelationId, nil)
		if err != nil && resCode != fiber.StatusCreated {
			return err
		}

		delete(inbox, coreelationId)

		return c.SendStatus(fiber.StatusCreated)
	}
}

func RollbackSaleOrderHandler(config SecondStepAppConfig, inbox CreateSaleOrderInbox) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var coreelationId = c.Params("correlationId")

		resCode, _, err := fasthttp.Post(nil, config.FirstStepAppURL+"/products/rollback/"+coreelationId, nil)
		if err != nil && resCode != fiber.StatusNoContent {
			return err
		}

		delete(inbox, coreelationId)

		return c.SendStatus(fiber.StatusOK)
	}
}
