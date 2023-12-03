package choreography

import "github.com/gofiber/fiber/v2"

type (
	CreateProductRequestInbox map[string]*CreateProductRequest
	ProductDB                 map[string]*Product
)

// ProductService
func NewFirstStepApp(inbox CreateProductRequestInbox, db ProductDB) *fiber.App {
	app := fiber.New()

	app.Post("/products", CreateProudctHandler(inbox))
	app.Post("/products/rollback/:correlationId", RollbackProductHandler(inbox))
	app.Post("/products/commit/:correlationId", CommitProductHandler(inbox, db))

	return app
}

type Product struct {
	Name string
}

type CreateProductRequest struct {
	CorrelationId string `json:"correlationId"`
	Name          string `json:"name"`
}

func CreateProudctHandler(inbox CreateProductRequestInbox) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateProductRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}

		inbox[req.CorrelationId] = &req

		return c.SendStatus(fiber.StatusAccepted)
	}
}

func CommitProductHandler(inbox CreateProductRequestInbox, db ProductDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var coreelationId = c.Params("correlationId")
		var awaitingRequest = inbox[coreelationId]

		db[coreelationId] = &Product{
			Name: awaitingRequest.Name,
		}

		delete(inbox, coreelationId)

		return c.SendStatus(fiber.StatusCreated)
	}
}

func RollbackProductHandler(inbox CreateProductRequestInbox) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var coreelationId = c.Params("correlationId")
		delete(inbox, coreelationId)

		return c.SendStatus(fiber.StatusNoContent)
	}
}
