package finnhub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

// ProfileParams represents parameters for company profile request
type ProfileParams struct {
	Symbol string // Stock symbol (e.g., "AAPL")
}

// ProfileResult represents the company profile data
type ProfileResult struct {
	Symbol    string  // Stock symbol
	Name      string  // Company name
	Country   string  // Country
	Exchange  string  // Exchange
	Industry  string  // Industry (finnhubIndustry)
	IPO       string  // IPO date
	MarketCap float64 // Market capitalization
	Phone     string  // Phone number
	WebURL    string  // Website URL
	Logo      string  // Logo URL
	Ticker    string  // Ticker
}

// GetProfile2 retrieves company profile data
func (c *Client) GetProfile2(ctx context.Context, params *ProfileParams) (*ProfileResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s?symbol=%s&token=%s", BaseURL, Profile2API, params.Symbol, c.apiKey)

	req := request.Request{
		Method:  "GET",
		URL:     url,
		Headers: DefaultHeaders,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseProfileResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type finnhubProfileResponse struct {
	Country   string  `json:"country"`
	Exchange  string  `json:"exchange"`
	Industry  string  `json:"finnhubIndustry"`
	IPO       string  `json:"ipo"`
	Logo      string  `json:"logo"`
	MarketCap float64 `json:"marketCapitalization"`
	Name      string  `json:"name"`
	Phone     string  `json:"phone"`
	Ticker    string  `json:"ticker"`
	WebURL    string  `json:"weburl"`
}

func parseProfileResponse(body []byte) (*ProfileResult, error) {
	var resp finnhubProfileResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &ProfileResult{
		Symbol:    resp.Ticker,
		Name:      resp.Name,
		Country:   resp.Country,
		Exchange:  resp.Exchange,
		Industry:  resp.Industry,
		IPO:       resp.IPO,
		MarketCap: resp.MarketCap,
		Phone:     resp.Phone,
		WebURL:    resp.WebURL,
		Logo:      resp.Logo,
		Ticker:    resp.Ticker,
	}, nil
}
