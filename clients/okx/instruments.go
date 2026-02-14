package okx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/souloss/quantds/request"
)

// API Endpoint Constants
const (
	EndpointInstruments = "/api/v5/public/instruments"

	// Query Parameters
	ParamInstType = "instType"

	// Instrument Types
	InstTypeSpot    = "SPOT"
	InstTypeMargin  = "MARGIN"
	InstTypeSwap    = "SWAP"
	InstTypeFutures = "FUTURES"
	InstTypeOption  = "OPTION"
)

// InstrumentParams represents parameters for instruments request
type InstrumentParams struct {
	InstType string // Instrument type: SPOT, MARGIN, SWAP, FUTURES, OPTION
	InstId   string // Instrument ID (optional, for filtering)
}

// InstrumentResult represents the instruments result
type InstrumentResult struct {
	Instruments []InstrumentData
	Total       int
}

// InstrumentData represents a single instrument/trading pair
type InstrumentData struct {
	InstID     string `json:"instId"`     // Instrument ID, e.g., "BTC-USDT"
	InstType   string `json:"instType"`   // Instrument type: SPOT, MARGIN, SWAP, FUTURES, OPTION
	BaseCcy    string `json:"baseCcy"`    // Base currency, e.g., "BTC"
	QuoteCcy   string `json:"quoteCcy"`   // Quote currency, e.g., "USDT"
	SettleCcy  string `json:"settleCcy"`  // Settlement currency
	CtVal      string `json:"ctVal"`      // Contract value
	CtMult     string `json:"ctMult"`     // Contract multiplier
	CtValCcy   string `json:"ctValCcy"`   // Contract value currency
	OptType    string `json:"optType"`    // Option type: C (call), P (put)
	Stk        string `json:"stk"`        // Strike price
	ListTime   string `json:"listTime"`   // Listing time
	ExpTime    string `json:"expTime"`    // Expiration time
	Lever      string `json:"lever"`      // Max leverage
	TickSz     string `json:"tickSz"`     // Tick size
	LotSz      string `json:"lotSz"`      // Lot size
	MinSz      string `json:"minSz"`      // Minimum order size
	CtType     string `json:"ctType"`     // Contract type
	Alias      string `json:"alias"`      // Contract date alias
	State      string `json:"state"`      // Instrument state: live
	Uly        string `json:"uly"`        // Underlying instrument
	Category   string `json:"category"`   // Fee category
	BaseCcyLogo string `json:"baseCcyLogo"` // Base currency logo URL
	QuoteCcyLogo string `json:"quoteCcyLogo"` // Quote currency logo URL
}

// GetInstruments retrieves list of instruments
func (c *Client) GetInstruments(ctx context.Context, params *InstrumentParams) (*InstrumentResult, *request.Record, error) {
	if params == nil {
		params = &InstrumentParams{}
	}

	u, _ := url.Parse(c.BaseURL + EndpointInstruments)
	q := u.Query()
	if params.InstType != "" {
		q.Add(ParamInstType, params.InstType)
	}
	if params.InstId != "" {
		q.Add(ParamInstID, params.InstId)
	}

	req := request.Request{
		Method: "GET",
		URL:    u.String() + "?" + q.Encode(),
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result Response
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, record, err
	}

	if result.Code != "0" {
		return nil, record, fmt.Errorf("api error: %s (code: %s)", result.Msg, result.Code)
	}

	var instruments []InstrumentData
	if err := json.Unmarshal(result.Data, &instruments); err != nil {
		return nil, record, err
	}

	return &InstrumentResult{
		Instruments: instruments,
		Total:       len(instruments),
	}, record, nil
}

// GetSpotInstruments retrieves all spot instruments
func (c *Client) GetSpotInstruments(ctx context.Context) (*InstrumentResult, *request.Record, error) {
	return c.GetInstruments(ctx, &InstrumentParams{
		InstType: InstTypeSpot,
	})
}

// GetSwapInstruments retrieves all swap (perpetual) instruments
func (c *Client) GetSwapInstruments(ctx context.Context) (*InstrumentResult, *request.Record, error) {
	return c.GetInstruments(ctx, &InstrumentParams{
		InstType: InstTypeSwap,
	})
}

// GetFuturesInstruments retrieves all futures instruments
func (c *Client) GetFuturesInstruments(ctx context.Context) (*InstrumentResult, *request.Record, error) {
	return c.GetInstruments(ctx, &InstrumentParams{
		InstType: InstTypeFutures,
	})
}
