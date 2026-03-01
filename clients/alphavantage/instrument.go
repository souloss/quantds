package alphavantage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

type SearchParams struct {
	Keywords string
}

type SearchResult struct {
	Matches []SearchMatch
	Count   int
}

type SearchMatch struct {
	Symbol      string `json:"1. symbol"`
	Name        string `json:"2. name"`
	Type        string `json:"3. type"`
	Region      string `json:"4. region"`
	MarketOpen  string `json:"5. marketOpen"`
	MarketClose string `json:"6. marketClose"`
	Timezone    string `json:"7. timezone"`
	Currency    string `json:"8. currency"`
}

func (c *Client) SearchSymbol(ctx context.Context, params *SearchParams) (*SearchResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s?function=SYMBOL_SEARCH&keywords=%s&apikey=%s",
		BaseURL, QueryAPI, params.Keywords, c.apiKey)

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

	result, err := parseSearchResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type searchResponse struct {
	BestMatches []SearchMatch `json:"bestMatches"`
}

func parseSearchResponse(body []byte) (*SearchResult, error) {
	var resp searchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &SearchResult{
		Matches: resp.BestMatches,
		Count:   len(resp.BestMatches),
	}, nil
}
