package choreography

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFirstStepApp(t *testing.T) {
	inbox := make(map[string]*CreateProductRequest, 0)
	db := make(map[string]*Product, 0)

	app := NewFirstStepApp(inbox, db)

	// Step 1: Create product Request
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer([]byte(`{"correlationId": "test", "name": "product 1"}`)))
	req.Header.Set("Content-Type", "application/json")

	res, err := app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, res.StatusCode)

	require.Equal(t, 1, len(inbox))
	require.Equal(t, 0, len(db))

	// Step 2: Commit product
	req = httptest.NewRequest(http.MethodPost, "/products/commit/test", nil)
	res, err = app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	require.Equal(t, 0, len(inbox))
	require.Equal(t, 1, len(db))

	// Step 3: Rollback product
	req = httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer([]byte(`{"correlationId": "test_rollback", "name": "product 2"}`)))
	req.Header.Set("Content-Type", "application/json")

	res, err = app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusAccepted, res.StatusCode)

	require.Equal(t, 1, len(inbox))
	require.Equal(t, 1, len(db))

	req = httptest.NewRequest(http.MethodPost, "/products/rollback/test_rollback", nil)
	res, err = app.Test(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	require.Equal(t, 0, len(inbox))
	require.Equal(t, 1, len(db))
}
