package tushare

import (
	"context"
	"testing"
	"time"
)

// TestClient_GetConcept tests retrieving concept list
// API Rule: Requires Tushare token
func TestClient_GetConcept(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, record, err := client.GetConcept(ctx, &ConceptParams{
		Src: "ts",
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetConcept() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d concept rows", len(rows))
	if len(rows) == 0 {
		t.Log("Warning: No concepts returned")
		return
	}

	// Print first 5 concepts
	for i, r := range rows {
		if i >= 5 {
			break
		}
		t.Logf("Concept[%d]: code=%s, name=%s, src=%s", i, r.Code, r.Name, r.Src)
	}
}

// TestClient_GetConceptDetail tests retrieving stocks in a concept
func TestClient_GetConceptDetail(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// First get a concept code
	concepts, _, err := client.GetConcept(ctx, &ConceptParams{})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetConcept() error = %v", err)
	}

	if len(concepts) == 0 {
		t.Skip("No concepts available for testing")
	}

	conceptCode := concepts[0].Code
	t.Logf("Testing with concept: %s (%s)", conceptCode, concepts[0].Name)

	// Get detail
	rows, record, err := client.GetConceptDetail(ctx, &ConceptDetailParams{
		ID: conceptCode,
	})
	skipOnTokenError(t, err)
	if err != nil {
		t.Fatalf("GetConceptDetail() error = %v", err)
	}
	if record == nil {
		t.Fatal("record is nil")
	}

	t.Logf("Got %d stocks in concept %s", len(rows), conceptCode)
	for i, r := range rows {
		if i >= 5 {
			break
		}
		t.Logf("Stock[%d]: code=%s, name=%s, in_date=%s",
			i, r.TSCode, r.Name, r.InDate)
	}
}
