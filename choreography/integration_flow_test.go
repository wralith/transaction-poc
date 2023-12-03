package choreography

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestIntegrationHappyPath(t *testing.T) {
	var (
		firstStepInbox = make(map[string]*CreateProductRequest)
		firstStepDB    = make(map[string]*Product)

		secondStepInbox = make(map[string]*CreateSaleOrderRequest)
		secondStepDB    = make(map[string]*SaleOrder)

		thirdStepInbox = make(map[string]*CreateTargetAudiencePatternRequest)
		thirdStepDB    = make(map[string]*TargetAudiencePattern)
	)

	firstStepApp := NewFirstStepApp(firstStepInbox, firstStepDB)
	secondStepApp := NewSecondStepApp(secondStepInbox, secondStepDB, GetSecondStepAppConfig())
	thirdStepApp := NewThirdStepApp(thirdStepInbox, thirdStepDB, GetThirdStepAppConfig())

	go firstStepApp.Listen(":3000")
	go secondStepApp.Listen(":3001")
	go thirdStepApp.Listen(":3002")

	// Acquire request and response
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	// Set request method and content type for next requests
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")

	// Create product
	req.SetRequestURI("http://localhost:3000/products")
	req.SetBody([]byte(`{"correlationId": "test", "name": "test"}`))

	// Send request
	err := fasthttp.Do(req, res)
	require.NoError(t, err)

	// Check if CreateProductRequest is in inbox
	require.Equal(t, 202, res.StatusCode())
	require.Len(t, firstStepInbox, 1)

	// Create sale order
	req.SetRequestURI("http://localhost:3001/orders")
	req.SetBody([]byte(`{"correlationId": "test", "name": "test", "price": 1}`))

	err = fasthttp.Do(req, res)
	require.NoError(t, err)

	require.Equal(t, 202, res.StatusCode())
	require.Len(t, secondStepInbox, 1)

	// Create target audience pattern
	req.SetRequestURI("http://localhost:3002/targets")
	req.SetBody([]byte(`{"correlationId": "test", "name": "test", "pattern": "test"}`))

	err = fasthttp.Do(req, res)
	require.NoError(t, err)

	require.Equal(t, 202, res.StatusCode())
	require.Len(t, thirdStepInbox, 1)

	// Commit target audience pattern
	// Expect that the all steps will be committed

	req.SetRequestURI("http://localhost:3002/targets/commit/test")
	req.SetBody([]byte(nil))
	err = fasthttp.Do(req, res)

	require.NoError(t, err)
	require.Equal(t, 201, res.StatusCode())

	// Check if the product is created
	require.Len(t, firstStepDB, 1)
	require.Equal(t, "test", firstStepDB["test"].Name)

	// Check if the sale order is created
	require.Len(t, secondStepDB, 1)
	require.Equal(t, "test", secondStepDB["test"].Name)

	// Check if the target audience pattern is created
	require.Len(t, thirdStepDB, 1)
	require.Equal(t, "test", thirdStepDB["test"].Name)

	// Check if inboxes cleaned
	require.Len(t, firstStepInbox, 0)
	require.Len(t, secondStepInbox, 0)
	require.Len(t, thirdStepInbox, 0)
}
