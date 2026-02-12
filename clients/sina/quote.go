package sina

import (
	"context"

	"github.com/souloss/quantds/request"
)

// QuoteParams represents parameters for real-time quote request
// This is an alias to SpotParams for API consistency
type QuoteParams = SpotParams

// QuoteResult represents the real-time quote result
// This is an alias to SpotResult for API consistency
type QuoteResult = SpotResult

// QuoteData represents a single real-time market quote
// This is an alias to SpotQuote for API consistency
type QuoteData = SpotQuote

// GetQuotes retrieves real-time quotes for multiple stocks
// This is an alias to GetSpot for API consistency
func (c *Client) GetQuotes(ctx context.Context, params *QuoteParams) (*QuoteResult, *request.Record, error) {
	return c.GetSpot(ctx, params)
}
