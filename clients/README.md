# Clients Module Documentation

This directory contains client implementations for various financial data providers, covering multiple markets (Crypto, Stocks, Forex, etc.) and asset classes.

The goal is to provide a unified, type-safe, and reliable way to access external APIs for quantitative data analysis.

## 🌟 Core Philosophy

- **Multi-Asset Support**: The system is designed to fetch data for any tradable asset (Equities, Crypto, Forex, Commodities).
- **Type Safety**: All API interactions must use strongly-typed structs for requests and responses.
- **Test-Driven**: Every API endpoint wrapper must have a corresponding live integration test.
- **Isolation**: Each client is a self-contained package with minimal external dependencies.

## 📁 Directory Structure

Each data provider should have its own subdirectory.

```text
clients/
├── binance/          # Crypto Exchange (Binance)
├── okx/              # Crypto Exchange (OKX)
├── coingecko/        # Crypto Market Data Aggregator
├── yahoo/            # Global Stocks & Forex (Yahoo Finance)
├── alphavantage/     # US Stocks & Forex (Alpha Vantage) [Example]
└── cninfo/           # China A-shares Information [Example]
```

### File Organization within a Client Package

| File | Purpose |
|------|---------|
| `client.go` | `Client` struct definition, `NewClient` constructor, base configuration (URL, Auth). |
| `<resource>.go` | Implementation of specific API endpoints (e.g., `ticker.go`, `kline.go`, `financials.go`). |
| `<resource>_test.go` | **Live integration tests** for the corresponding resource file. |
| `types.go` | (Optional) Shared data structures if used across multiple files. |

## 🛠 Implementation Guidelines

### 1. API Implementation (One File Per Resource)
Do not dump all methods into a single file. Group methods by resource or functionality.

- **Bad**: `api.go` containing 20 different methods.
- **Good**: `quote.go` (Real-time price), `historical.go` (History), `search.go` (Symbol lookup).

### 2. Type Safety & Constants
Avoid "magic strings" and untyped maps (`map[string]interface{}`).

- **Constants**: Define API endpoints and query parameters as constants at the top of the file.
- **Request Structs**: Create dedicated structs for input parameters.
- **Response Structs**: Create dedicated structs matching the JSON response.

**Example:**
```go
// ticker.go
const (
    EndpointTicker = "/api/v5/market/ticker"
    ParamInstID    = "instId"
)

type TickerRequest struct {
    InstID string // e.g., "BTC-USDT" or "AAPL"
}

type TickerResponse struct {
    LastPrice string `json:"last"`
    Volume    string `json:"vol"`
}

func (c *Client) GetTicker(ctx context.Context, req *TickerRequest) (*TickerResponse, error) { ... }
```

### 3. Testing Requirements (Strict)
Every implemented API function **MUST** have at least one corresponding test case in a `_test.go` file.

- **Live Data**: Tests should connect to the real API (unless authentication costs money or is impossible).
- **Verification**: Tests must verify that data is actually returned (e.g., `len(data) > 0`).
- **Logging**: Use `t.Logf` to print sample data (e.g., latest price) to prove the integration works.
- **Geo-Restriction Handling**: If an API is blocked in certain regions (e.g., Mainland China, USA), use helper functions to skip tests gracefully instead of failing.

**Example Test:**
```go
// ticker_test.go
func TestClient_GetTicker(t *testing.T) {
    client := NewClient(nil)
    
    // 1. Call API
    resp, _, err := client.GetTicker(context.Background(), &TickerRequest{InstID: "AAPL"})
    
    // 2. Handle known network restrictions
    if err != nil {
        checkAPIError(t, err) // Helper to skip if 403/451
        return
    }

    // 3. Log & Verify
    t.Logf("Latest Price: %s", resp.LastPrice)
    if resp.LastPrice == "" {
        t.Error("Expected price, got empty string")
    }
}
```

## 🌍 Supported Markets & Providers

This system is extensible. When adding a new provider, consider:

1.  **Crypto**: Binance, OKX, CoinGecko, etc.
2.  **US Stocks**: Yahoo Finance, Alpha Vantage, Polygon.io.
3.  **CN Stocks (A-Share)**: Sina Finance, EastMoney, Tushare.
4.  **Forex/Commodities**: OANDA, Fixer.io.

## 📝 Contribution Checklist

Before submitting a new client or feature:
- [ ] Created dedicated file for the resource (e.g., `dividends.go`).
- [ ] Defined strong types for Request and Response.
- [ ] Added `_test.go` file with a live test case.
- [ ] Test prints sample data (`t.Logf`) and asserts validity.
- [ ] All tests pass (`go test -v ./clients/<provider>/...`).
