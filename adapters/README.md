# Adapters Module Documentation

This directory contains adapter implementations that bridge low-level **clients** (API wrappers) with high-level **domain** interfaces. Adapters implement the `manager.Provider` interface, enabling the system to use multiple data sources interchangeably through a unified abstraction.

## Architecture Overview

```text
domain/ (interfaces)  ←  adapters/ (bridge)  ←  clients/ (API wrappers)
     ↑                      ↑                       ↑
 spot.Source            manager.Provider         binance.Client
 kline.Source           (Name, Fetch,            okx.Client
 instrument.Source       SupportedMarkets,       eastmoney.Client
 ...                     CanHandle)              ...
```

Each adapter:
1. Takes a **client** instance as constructor parameter
2. Implements `manager.Provider[Req, Resp]` interface
3. Converts between **domain types** and **client types**
4. Tracks request metadata via `manager.RequestTrace`

---

## Provider Interface

All adapters implement the generic `manager.Provider` interface:

```go
type Provider[Req, Resp any] interface {
    Name() string
    Fetch(ctx context.Context, client request.Client, req Req) (Resp, *RequestTrace, error)
    SupportedMarkets() []domain.Market
    CanHandle(symbol string) bool
}
```

### Interface Methods

| Method | Purpose |
|--------|---------|
| `Name()` | Returns a constant identifier (e.g., `"eastmoney"`, `"binance"`) |
| `Fetch()` | Converts domain request → client call → domain response |
| `SupportedMarkets()` | Returns list of supported markets (e.g., `[MarketCN]`, `[MarketCrypto]`) |
| `CanHandle(symbol)` | Checks if the adapter can handle a given symbol format |

---

## Directory Structure

```text
adapters/
├── binance/       # Crypto: K线, 行情, 证券列表
├── okx/           # Crypto: K线, 行情, 证券列表
├── yahoo/         # US: K线, 行情, 证券列表
├── sina/          # CN: K线, 行情
├── tencent/       # CN: K线, 行情, 行情(Quote)
├── eastmoney/     # CN: K线, 行情, 证券列表, 财务, 公告, 个股档案
├── eastmoneyhk/   # HK: K线, 行情, 证券列表
├── tushare/       # CN: K线, 行情, 证券列表, 财务, 公告, 个股档案
├── xueqiu/        # CN: K线, 行情, 证券列表, 个股档案
├── cninfo/        # CN: 证券列表, 公告
├── sse/           # CN: 证券列表 (上交所)
├── szse/          # CN: 证券列表 (深交所)
└── bse/           # CN: 证券列表 (北交所)
```

### File Organization within an Adapter Package

| File | Purpose |
|------|---------|
| `kline.go` | K 线适配器 — 实现 `manager.Provider[kline.Request, kline.Response]` |
| `spot.go` | 实时行情适配器 — 实现 `manager.Provider[spot.Request, spot.Response]` |
| `instrument.go` | 证券列表适配器 — 实现 `manager.Provider[instrument.Request, instrument.Response]` |
| `financial.go` | 财务数据适配器 — 实现 `manager.Provider[financial.Request, financial.Response]` |
| `announcement.go` | 公告新闻适配器 — 实现 `manager.Provider[announcement.Request, announcement.Response]` |
| `profile.go` | 个股档案适配器 — 实现 `manager.Provider[profile.Request, profile.Response]` |
| `*_test.go` | 每个适配器的单元测试 |

---

## Supported Markets & Providers

| Market | K线 | 行情 | 证券列表 | 财务 | 公告 | 个股档案 |
|--------|-----|------|----------|------|------|----------|
| **CN (A股)** | eastmoney, sina, tencent, tushare, xueqiu | sina, tencent, eastmoney, xueqiu | eastmoney, tushare, cninfo, sse, szse, bse | eastmoney, tushare | eastmoney, cninfo | eastmoney, tushare, xueqiu |
| **HK (港股)** | eastmoneyhk | eastmoneyhk | eastmoneyhk | - | - | - |
| **US (美股)** | yahoo | yahoo | yahoo | - | - | - |
| **Crypto** | binance, okx | binance, okx | binance, okx | - | - | - |

---

## Implementation Guidelines

### 1. Adapter Structure

Every adapter must follow the same structural pattern:

```go
package <provider>

import (
    "github.com/souloss/quantds/clients/<provider>"
    "github.com/souloss/quantds/domain"
    "github.com/souloss/quantds/domain/kline"
    "github.com/souloss/quantds/manager"
    "github.com/souloss/quantds/request"
)

// Name constant — defined once per package, shared by all adapters
const Name = "<provider>"

// Supported markets — defined once per package, shared by all adapters
var supportedMarkets = []domain.Market{domain.MarketXX}

// XxxAdapter implements manager.Provider for a specific domain
type XxxAdapter struct {
    client *<provider>.Client
}

func NewXxxAdapter(client *<provider>.Client) *XxxAdapter {
    return &XxxAdapter{client: client}
}

func (a *XxxAdapter) Name() string                      { return Name }
func (a *XxxAdapter) SupportedMarkets() []domain.Market { return supportedMarkets }
func (a *XxxAdapter) CanHandle(symbol string) bool       { /* ... */ }
func (a *XxxAdapter) Fetch(ctx context.Context, _ request.Client, req Req) (Resp, *manager.RequestTrace, error) {
    // 1. Create trace
    trace := manager.NewRequestTrace(Name)
    
    // 2. Convert domain request to client params
    // 3. Call client method
    // 4. trace.AddRequest(record)
    // 5. Convert client response to domain response
    // 6. trace.Finish()
    // 7. Return
}

// Compile-time interface check
var _ manager.Provider[Req, Resp] = (*XxxAdapter)(nil)
```

### 2. Name & SupportedMarkets Constants

- `Name` must be a **package-level constant** (not hardcoded string in method)
- `supportedMarkets` must be a **package-level variable** shared across all adapters in the same package
- Both are defined **once** in the first adapter file (typically `kline.go`) and reused

### 3. No Magic Strings

- Use `domain.ExchangeBinance` instead of `"BINANCE"`
- Use `domain.ExchangeNYSE` instead of `"NYSE"`
- Use `domain.MarketHK` instead of `"HK"`
- Use `string(domain.ExchangeXxx)` for type conversions when comparing with client API strings

### 4. Symbol Format Handling

Adapters must handle symbol format conversion between domain and client:

- **Domain format**: `000001.SZ`, `600519.SH`, `BTCUSDT`, `AAPL.US`, `00700.HK.HKEX`
- **Client format**: Provider-specific (e.g., `sz000001` for Sina, `BTC-USDT` for OKX)

Each adapter should use the client's conversion functions (e.g., `client.ToXxxSymbol()`) and handle edge cases:

```go
func (a *SpotAdapter) Fetch(...) {
    // Convert domain symbols to client-specific format
    for _, s := range req.Symbols {
        symbol, err := client.ToXxxSymbol(s)
        if err != nil {
            continue // Skip invalid symbols gracefully
        }
        symbols = append(symbols, symbol)
    }
}
```

### 5. RequestTrace Usage

Every `Fetch` method must:
1. Create a trace: `trace := manager.NewRequestTrace(Name)`
2. Record each HTTP request: `trace.AddRequest(record)`
3. Finish the trace: `trace.Finish()`
4. Return the trace alongside the response

### 6. Error Handling

- **Critical errors** (main data source fails): Return `(empty, trace, err)`
- **Supplementary data errors** (e.g., BS/CF in financial adapter): Log but continue with partial data
- **Symbol conversion errors** in batch operations: Skip invalid symbols, continue with valid ones

### 7. Compile-Time Interface Check

Every adapter file must include a compile-time check at the bottom:

```go
var _ manager.Provider[kline.Request, kline.Response] = (*KlineAdapter)(nil)
```

---

## Testing Guidelines

### Unit Tests (Required)

Each adapter must have unit tests that verify:

1. **Constructor**: `NewXxxAdapter` returns non-nil
2. **Name**: Returns the correct constant
3. **SupportedMarkets**: Returns the expected market list
4. **CanHandle**: Correctly identifies supported and unsupported symbols

```go
func TestNewKlineAdapter(t *testing.T) {
    client := <provider>.NewClient()
    adapter := NewKlineAdapter(client)
    
    if adapter == nil {
        t.Error("NewKlineAdapter returned nil")
    }
    if adapter.Name() != Name {
        t.Errorf("Expected name '%s', got '%s'", Name, adapter.Name())
    }
}

func TestKlineAdapter_CanHandle(t *testing.T) {
    client := <provider>.NewClient()
    adapter := NewKlineAdapter(client)
    
    tests := []struct {
        symbol    string
        canHandle bool
    }{
        {"000001.SZ", true},   // positive case
        {"BTCUSDT", false},    // negative case
    }
    for _, tt := range tests {
        t.Run(tt.symbol, func(t *testing.T) {
            if got := adapter.CanHandle(tt.symbol); got != tt.canHandle {
                t.Errorf("CanHandle(%s) = %v, want %v", tt.symbol, got, tt.canHandle)
            }
        })
    }
}
```

### Important Test Notes

- **Do NOT pass `nil`** to `NewClient()`. Use `NewClient()` (no arguments) instead — `nil` causes a nil pointer dereference panic with variadic option constructors.
- **HK symbols**: Use `"00700.HK.HKEX"` or `"00700.HKEX"` format (not `"00700.HK"`, as `.HK` is not recognized as `HKEX` exchange by the domain parser).

---

## Adding a New Adapter

### Checklist

- [ ] Identified the target domain (kline, spot, instrument, financial, announcement, profile)
- [ ] Confirmed corresponding client methods exist in `clients/<provider>/`
- [ ] Created adapter file implementing `manager.Provider` interface
- [ ] Used package-level `Name` constant and `supportedMarkets` variable
- [ ] Added compile-time interface check (`var _ manager.Provider[...] = (...)`)
- [ ] No magic strings — used `domain.Exchange*`, `domain.Market*` constants
- [ ] Added unit test file testing constructor, Name, SupportedMarkets, CanHandle
- [ ] Tests pass: `go test ./adapters/<provider>/...`
- [ ] Vet passes: `go vet ./adapters/<provider>/...`
- [ ] Registered adapter in `facade/service.go` with appropriate priority
- [ ] Updated this README's Supported Markets table

### Priority Guidelines (for Facade Registration)

| Priority | Value | Use Case |
|----------|-------|----------|
| `PriorityHighest` (100) | Primary/most reliable data source |
| `PriorityHigh` (75) | Strong backup source |
| `PriorityMedium` (50) | Decent fallback |
| `PriorityLow` (25) | Low-priority fallback |
| `PriorityLowest` (1) | Last resort |
