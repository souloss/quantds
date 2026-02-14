package eastmoney

import (
	"context"
	"testing"
)

// TestClient_GetConceptList tests retrieving concept list
// API Rule: No authentication required
// Geo-Restriction: May be blocked in some regions
func TestClient_GetConceptList(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	params := &ConceptListParams{
		PageSize: 10,
		PageNo:   1,
	}

	result, record, err := client.GetConceptList(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Concept List Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d concepts", len(result))

	if len(result) == 0 {
		t.Log("Warning: No concepts returned")
		return
	}

	for i, concept := range result {
		t.Logf("Concept[%d]: code=%s, name=%s, change=%.2f%%, netInflow=%.2f",
			i, concept.Code, concept.Name, concept.ChangePercent, concept.MainNetInflow)
	}
}

// TestClient_GetConceptStocks tests retrieving stocks in a concept
func TestClient_GetConceptStocks(t *testing.T) {
	client := NewClient()
	defer client.Close()
	ctx := context.Background()

	// First get a concept code
	conceptList, _, err := client.GetConceptList(ctx, &ConceptListParams{PageSize: 1})
	if err != nil {
		checkAPIError(t, err)
		return
	}

	if len(conceptList) == 0 {
		t.Skip("No concepts available for testing")
	}

	conceptCode := conceptList[0].Code
	t.Logf("Testing with concept: %s (%s)", conceptCode, conceptList[0].Name)

	// Get stocks in the concept
	params := &ConceptStocksParams{
		ConceptCode: conceptCode,
		PageSize:    10,
		PageNo:      1,
	}

	result, record, err := client.GetConceptStocks(ctx, params)
	if err != nil {
		checkAPIError(t, err)
		return
	}

	t.Logf("Concept Stocks Response Status: %d", record.Response.StatusCode)
	t.Logf("Got %d stocks in concept %s", len(result), conceptCode)

	for i, stock := range result {
		t.Logf("Stock[%d]: code=%s, name=%s, price=%.2f, change=%.2f%%",
			i, stock.Code, stock.Name, stock.Latest, stock.ChangeRate)
	}
}
