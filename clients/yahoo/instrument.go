package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/souloss/quantds/request"
)

// InstrumentParams represents parameters for instrument list request
type InstrumentParams struct {
	Query    string // Search query (symbol or company name)
	Exchange string // Exchange filter (NASDAQ, NYSE, AMEX)
	Limit    int    // Number of results
}

// InstrumentResult represents the instrument list result
type InstrumentResult struct {
	Instruments []InstrumentData
	Total       int
}

// InstrumentData represents a single instrument
type InstrumentData struct {
	Symbol      string // Stock symbol
	Name        string // Company name
	Exchange    string // Exchange (NASDAQ, NYSE, AMEX)
	AssetType   string // Asset type (Stock, ETF, etc.)
	Currency    string // Currency (USD)
	ListingDate string // Listing date
}

// GetInstruments retrieves US stock instruments by search query
// Uses Yahoo Finance search API
func (c *Client) GetInstruments(ctx context.Context, params *InstrumentParams) (*InstrumentResult, *request.Record, error) {
	if params.Limit <= 0 {
		params.Limit = 30
	}

	// Build search query
	query := params.Query
	if query == "" {
		// Default to common stocks
		query = ""
	}

	// Use Yahoo Finance search API
	searchURL := fmt.Sprintf("%s%s?q=%s&newsCount=0&listsCount=0&quotesCount=1&sort=SIMILARITY",
		BaseURL, SearchAPI, url.QueryEscape(query))

	req := request.Request{
		Method:  "GET",
		URL:     searchURL,
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

// GetAllUSStocks retrieves all major US stocks
// This uses predefined lists of major US exchanges
func (c *Client) GetAllUSStocks(ctx context.Context) (*InstrumentResult, *request.Record, error) {
	// Return a list of major US stock symbols
	// In production, this would fetch from a more comprehensive source
	majorStocks := []InstrumentData{
		// NASDAQ top stocks
		{Symbol: "AAPL", Name: "Apple Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "MSFT", Name: "Microsoft Corporation", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "GOOGL", Name: "Alphabet Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "AMZN", Name: "Amazon.com Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "NVDA", Name: "NVIDIA Corporation", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "META", Name: "Meta Platforms Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "TSLA", Name: "Tesla Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "BRK.B", Name: "Berkshire Hathaway Inc.", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "LLY", Name: "Eli Lilly and Company", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "AVGO", Name: "Broadcom Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "V", Name: "Visa Inc.", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "JPM", Name: "JPMorgan Chase & Co.", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "WMT", Name: "Walmart Inc.", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "XOM", Name: "Exxon Mobil Corporation", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "UNH", Name: "UnitedHealth Group Inc.", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "MA", Name: "Mastercard Incorporated", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "PG", Name: "Procter & Gamble Company", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "JNJ", Name: "Johnson & Johnson", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "HD", Name: "The Home Depot Inc.", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "MRK", Name: "Merck & Co. Inc.", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "COST", Name: "Costco Wholesale Corporation", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "ABBV", Name: "AbbVie Inc.", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "AMD", Name: "Advanced Micro Devices Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "ADBE", Name: "Adobe Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "CRM", Name: "Salesforce Inc.", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "NFLX", Name: "Netflix Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "CSCO", Name: "Cisco Systems Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "PEP", Name: "PepsiCo Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "KO", Name: "The Coca-Cola Company", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "INTC", Name: "Intel Corporation", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		// More popular ETFs
		{Symbol: "SPY", Name: "SPDR S&P 500 ETF Trust", Exchange: "NYSE", AssetType: "ETF", Currency: "USD"},
		{Symbol: "QQQ", Name: "Invesco QQQ Trust", Exchange: "NASDAQ", AssetType: "ETF", Currency: "USD"},
		{Symbol: "IWM", Name: "iShares Russell 2000 ETF", Exchange: "NYSE", AssetType: "ETF", Currency: "USD"},
		{Symbol: "VTI", Name: "Vanguard Total Stock Market ETF", Exchange: "NYSE", AssetType: "ETF", Currency: "USD"},
		{Symbol: "VOO", Name: "Vanguard S&P 500 ETF", Exchange: "NYSE", AssetType: "ETF", Currency: "USD"},
		// Popular crypto-related stocks
		{Symbol: "COIN", Name: "Coinbase Global Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "MARA", Name: "Marathon Digital Holdings Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "RIOT", Name: "Riot Platforms Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		// Chinese ADRs
		{Symbol: "BABA", Name: "Alibaba Group Holding Limited", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "JD", Name: "JD.com Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "PDD", Name: "PDD Holdings Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "NIO", Name: "NIO Inc.", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
		{Symbol: "BILI", Name: "Bilibili Inc.", Exchange: "NASDAQ", AssetType: "Stock", Currency: "USD"},
		{Symbol: "TAL", Name: "Tal Education Group", Exchange: "NYSE", AssetType: "Stock", Currency: "USD"},
	}

	return &InstrumentResult{
		Instruments: majorStocks,
		Total:       len(majorStocks),
	}, &request.Record{}, nil
}

type searchResponse struct {
	Quotes []searchQuote `json:"quotes"`
}

type searchQuote struct {
	Symbol    string `json:"symbol"`
	Shortname string `json:"shortname"`
	Longname  string `json:"longname"`
	ExchDisp  string `json:"exchDisp"`
	QuoteType string `json:"quoteType"`
	Currency  string `json:"currency"`
}

func parseSearchResponse(body []byte) (*InstrumentResult, error) {
	var resp searchResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	instruments := make([]InstrumentData, 0, len(resp.Quotes))
	for _, q := range resp.Quotes {
		// Filter to only stocks and ETFs
		if q.QuoteType != "EQUITY" && q.QuoteType != "ETF" {
			continue
		}
		// Skip if symbol contains special characters (like warrants)
		if strings.Contains(q.Symbol, "-") || strings.Contains(q.Symbol, ".") {
			continue
		}

		name := q.Shortname
		if name == "" {
			name = q.Longname
		}

		instruments = append(instruments, InstrumentData{
			Symbol:    q.Symbol,
			Name:      name,
			Exchange:  q.ExchDisp,
			AssetType: q.QuoteType,
			Currency:  q.Currency,
		})
	}

	return &InstrumentResult{
		Instruments: instruments,
		Total:       len(instruments),
	}, nil
}

// GetInstrumentsByExchange retrieves instruments filtered by exchange
// Uses the exchange name as search query to find relevant instruments
func (c *Client) GetInstrumentsByExchange(ctx context.Context, exchange string) (*InstrumentResult, *request.Record, error) {
	params := &InstrumentParams{
		Query:    exchange, // Use exchange name as search query
		Exchange: exchange,
		Limit:    100,
	}
	return c.GetInstruments(ctx, params)
}
