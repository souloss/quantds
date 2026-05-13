# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

quantds is a unified financial data service library (Go module: `github.com/souloss/quantds`) providing standardized access to multi-market (A-share, HK, US, Crypto, Forex) and multi-asset (stock, fund, bond, crypto) data. It wraps multiple data sources with automatic failover, priority-based provider selection, and two-level caching.

## Build & Test Commands

```bash
# Build
go build ./...

# Run all tests
go test -v ./...

# Run tests for a specific package
go test -v ./adapters/eastmoney/
go test -v ./clients/eastmoney/
go test -v ./manager/
go test -v ./facade/

# Run a single test
go test -v -run TestKlineAdapter_CanHandle ./adapters/eastmoney/

# Format code
go fmt ./...

# Lint (requires golangci-lint)
golangci-lint run

# Generate supported data sources table in README
go run cmd/gendoc/main.go

# Clean
go clean && rm -rf dist/
```

Facade tests (`facade/service_test.go`) call real external APIs. They use `checkFacadeError` to skip on network/API failures rather than failing the build, so they are safe to run in CI.

## Architecture: Facade - Manager - Adapter - Client

The project follows a strict 4-layer architecture. Data flows from top to bottom:

### 1. Facade (`facade/`)
`facade.Service` is the sole external entry point. It holds per-market `Manager` instances for each domain (kline, spot, instrument, profile, financial, announcement) and routes requests by parsing the symbol to determine the market. All provider registration and priority assignment happens in `initManagers()`.

### 2. Manager (`manager/`)
`Manager[Req, Resp]` is a generic core scheduler. Key responsibilities:
- **Provider selection**: Uses `Selector` (default: `PrioritySelector`) to order providers by priority
- **Failover**: Iterates providers in priority order; if the highest-priority fails, tries the next
- **Two-level caching**: `TwoLevelCache` with separate request-level and fetch-level TTLs
- **Metrics**: Pluggable `Collector` interface (`NoopCollector` default, `MemoryCollector` available)
- **Tracing**: `RequestTrace` records HTTP request/response details per provider attempt

The `Provider[Req, Resp]` interface requires: `Name()`, `Fetch(ctx, client, req)`, `SupportedMarkets()`, `CanHandle(symbol)`.

### 3. Adapter (`adapters/<datasource>/`)
Each adapter implements `manager.Provider[Req, Resp]` for a specific domain (kline, spot, instrument, etc.). Adapters:
- Take a client instance via constructor (e.g., `NewKlineAdapter(client)`)
- Convert domain request types to client-specific API parameters
- Convert client response data to domain response types
- Use `manager.NewRequestTrace()` to record HTTP interactions
- Compile-time interface check via `var _ manager.Provider[...] = (*Adapter)(nil)`

One adapter file per domain capability (e.g., `kline.go`, `spot.go`, `instrument.go`). Adapter packages are named by datasource (e.g., `adapters/eastmoney`, `adapters/binance`).

### 4. Client (`clients/<datasource>/`)
Clients are pure HTTP API wrappers with no business logic. They:
- Handle request construction, URL building, parameter encoding
- Parse raw API responses into structured Go types
- Return `(Result, *request.Record, error)` tuples for tracing
- Accept `request.Client` via `WithHTTPClient()` option for shared HTTP infrastructure
- Some clients require API keys via env vars (e.g., `ALPHAVANTAGE_API_KEY`, `FINNHUB_API_KEY`)

### 5. Domain (`domain/`)
Defines canonical Request/Response types for each data domain:
- `domain/kline` - K-line (candlestick) data: `Request{Symbol, Timeframe, StartTime, EndTime, Adjust}`, `Response{Symbol, Bars, Source}`, `Bar{Timestamp, OHLCV, Change, ChangeRate, TurnoverRate}`
- `domain/spot` - Real-time quotes
- `domain/instrument` - Security listings
- `domain/profile` - Security profiles
- `domain/financial` - Financial reports
- `domain/announcement` - News/announcements
- `domain/symbol.go` - `Symbol` struct with `Parse()` supporting formats: `CODE.EXCHANGE` (e.g., `000001.SZ`), `CODE.MARKET.EXCHANGE` (e.g., `00700.HK.HKEX`), prefix format (`SH600001`), and smart auto-detection

### 6. Request (`request/`)
HTTP client abstraction built on `resty.dev/v3` + `failsafe-go`. Provides:
- `Client` interface with `Do(ctx, Request) -> (Response, *Record, error)`
- Resilience policies: timeout, circuit breaker, retry, rate limiter (via failsafe-go)
- `NoopClient` and `CachingClient` for testing

## Adding a New Data Source

To add a new provider (e.g., "newsource"):
1. Create `clients/newsource/` with `client.go` (constructor + options pattern) and domain-specific files (e.g., `kline.go`, `spot.go`)
2. Create `adapters/newsource/` with adapter files per domain capability (e.g., `kline.go` implementing `manager.Provider[kline.Request, kline.Response]`)
3. Register adapters in `facade/service.go` `initManagers()` with appropriate priority and market
4. Update `cmd/gendoc/main.go` `providerNames` and `providerMarkets` maps
5. Run `go run cmd/gendoc/main.go` to update README tables

## Key Conventions

- Module path: `github.com/souloss/quantds`
- Go version: 1.24.0
- Markets: `CN`, `HK`, `US`, `CRYPTO`, `FOREX` (defined in `domain/symbol.go`)
- Priority constants in facade: `PriorityHighest=100`, `PriorityHigh=75`, `PriorityMedium=50`, `PriorityLow=25`, `PriorityLowest=1`
- Cache TTLs: Kline=5min, Spot=10sec, List=1hour
- Clients that need API keys read from env vars at `NewClient()` time
- Adapter compile-time checks: `var _ manager.Provider[Req, Resp] = (*Adapter)(nil)`
- Tests use standard `testing` package; client tests often have `testing_utils_test.go` helpers