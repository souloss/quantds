# Clients Module Documentation

This directory contains client implementations for various financial data providers, covering multiple markets (Crypto, Stocks, Forex, etc.) and asset classes.

## 🎯 Primary Goal: Domain Alignment

The **primary goal** of this module is to provide data that fulfills the requirements defined in the `domain` module. While clients are low-level API wrappers, their implementation should be guided by the needs of the core business logic.

### Guidelines for Domain Alignment

1.  **Consult Domain Definitions**: Before implementing a client, always review the `domain/` directory to understand what data models and interfaces are required by the system.
2.  **Prioritize Domain Needs**: When choosing which API endpoints to implement first, prioritize those that map directly to core domain concepts (e.g., K-line history, real-time quotes, asset lists).
3.  **Data Compatibility**: While client-specific request/response structs should match the external API exactly, keep the domain models in mind. Ensure that the data retrieved can be losslessly or easily converted to the corresponding domain types.

---

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
├── sina/             # CN Stocks — Sina Finance (新浪财经)
├── tencent/          # CN Stocks — Tencent Finance (腾讯证券)
├── eastmoney/        # CN Stocks — EastMoney (东方财富)
├── eastmoneyhk/      # HK Stocks — EastMoney HK (东方财富港股)
├── tushare/          # CN Stocks — Tushare (挖地兔)
├── xueqiu/           # CN Stocks — Xueqiu (雪球)
├── cninfo/           # CN Stocks — CNInfo (巨潮资讯)
├── sse/              # CN Exchange — Shanghai Stock Exchange (上交所)
├── szse/             # CN Exchange — Shenzhen Stock Exchange (深交所)
└── bse/              # CN Exchange — Beijing Stock Exchange (北交所)
```

### File Organization within a Client Package

| File | Purpose |
|------|---------|
| `client.go` | `Client` struct definition, `NewClient` constructor, base URL, default headers, helper methods. |
| `<resource>.go` | **Raw API Implementation**: Type-safe wrappers for specific API endpoints (e.g., `ticker.go`, `candles.go`). |
| `<resource>_test.go` | **Live integration tests** for the corresponding resource file. |
| `testing_utils_test.go` | Shared test helpers (e.g., `checkAPIError`) for graceful error handling. |
| `util.go` | (Optional) Shared helper functions (e.g., symbol parsing, type conversions). |
| `types.go` | (Optional) Shared data structures if used across multiple files. |

## 🛠 Implementation Guidelines

### 1. API Implementation (One File Per Resource)
Do not dump all methods into a single file. Group methods by resource or functionality.

- **Bad**: `api.go` containing 20 different methods.
- **Good**: `ticker.go` (Real-time price), `candles.go` (History), `coins_list.go` (Symbol lookup).

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
    client := NewClient()
    
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

### 4. Header & URL Constants
Avoid hardcoding HTTP headers and base URLs in resource files.

- **Headers**: Define `DefaultUserAgent`, `DefaultReferer`, etc. as constants in `client.go`.
- **BaseURL**: Define `BaseURL` as a constant in `client.go`. All resource files should reference it instead of hardcoding the URL.
- **Helper Methods**: For clients with authentication (cookies/tokens), create a `buildHeaders()` helper method to reduce duplication.

**Example:**
```go
// client.go
const (
    BaseURL          = "https://api.example.com"
    DefaultUserAgent = "Mozilla/5.0 ..."
    DefaultReferer   = "https://example.com/"
)

func (c *Client) buildHeaders() map[string]string {
    headers := map[string]string{
        "User-Agent": DefaultUserAgent,
        "Referer":    DefaultReferer,
    }
    if c.token != "" {
        headers["Authorization"] = "Bearer " + c.token
    }
    return headers
}
```

### 5. Testing Utilities (`testing_utils_test.go`)
Each client package should include a `testing_utils_test.go` file with a shared `checkAPIError` helper function.

- **Purpose**: Gracefully skip tests when APIs are unavailable due to geo-restrictions, authentication requirements, rate limiting, or service outages.
- **Status Codes to Handle**: 401 (Unauthorized), 403 (Forbidden), 429 (Rate Limited), 451 (Geo-blocked), 503 (Service Unavailable).
- **API-Level Errors**: Also handle non-HTTP errors like "Service not found", parse failures, and "client error" responses.
- **Pattern**: Always `t.Skipf` (not `t.Fatalf`) for expected/environmental failures.

**Example:**
```go
// testing_utils_test.go
func checkAPIError(t *testing.T, err error) {
    if err == nil {
        return
    }
    var reqErr *request.RequestError
    if errors.As(err, &reqErr) {
        switch reqErr.StatusCode {
        case 401, 403, 429, 451, 503:
            t.Skipf("Skipping: API restriction (status %d): %v", reqErr.StatusCode, err)
        }
    }
    // Handle API-level errors
    errMsg := err.Error()
    if strings.Contains(errMsg, "client error") ||
        strings.Contains(errMsg, "unmarshal") {
        t.Skipf("Skipping: API error: %v", err)
    }
    t.Fatalf("API request failed: %v", err)
}
```

### 6. Symbol Format Convention
Use a consistent symbol format across all Chinese market clients: `{code}.{exchange}`.

- **Format**: `000001.SZ`, `600519.SH`, `430001.BJ`
- **Do not** use prefix format like `SZ000001` or `SH600519` in public APIs or tests.
- Each client should provide internal conversion functions (e.g., `toSinaSymbol`, `toXueqiuSymbol`) to convert from the standard format to the provider-specific format.

### 7. File Naming Rules
- Test files **MUST** follow the `_test.go` suffix convention (e.g., `ticker_test.go`).
- **Never** create files like `*_test_new.go` or `*_test_v2.go` — these are not recognized by the Go test framework and will cause compilation errors when referencing test-only symbols.
- The `util.go` or `utils.go` file is acceptable for shared helper functions within a package.

## 🌍 Supported Markets & Providers

This system is extensible. When adding a new provider, consider:

1.  **Crypto**: Binance, OKX, CoinGecko, etc.
2.  **US Stocks**: Yahoo Finance, Alpha Vantage, Polygon.io.
3.  **CN Stocks (A-Share)**: Sina Finance, EastMoney, Tencent, Tushare, Xueqiu.
4.  **HK Stocks**: EastMoney HK.
5.  **Exchange Data (CN)**: SSE (Shanghai), SZSE (Shenzhen), BSE (Beijing), CNInfo.
6.  **Forex/Commodities**: OANDA, Fixer.io.

## 📝 Contribution Checklist

Before submitting a new client or feature:
- [ ] Checked `domain/` directory to understand current data requirements.
- [ ] Created dedicated file for the resource (e.g., `dividends.go`).
- [ ] Defined strong types for Request and Response.
- [ ] All constants (endpoints, query params, headers) defined at file/package level — no magic strings.
- [ ] Used `BaseURL` constant and `buildHeaders()` helper — no hardcoded URLs or headers in resource files.
- [ ] Added `testing_utils_test.go` with `checkAPIError` helper.
- [ ] Added `_test.go` file with a live test case.
- [ ] Test prints sample data (`t.Logf`) and asserts validity.
- [ ] Test uses `checkAPIError` for all error paths to handle geo-restrictions gracefully.
- [ ] All tests pass (`go test -v ./clients/<provider>/...`).
- [ ] `go vet ./clients/<provider>/...` passes with no warnings.
