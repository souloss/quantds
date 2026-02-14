# Facade Module Documentation

The `facade` package provides a unified entry point (`Service`) for accessing financial data across multiple markets (A 股、港股、美股、加密货币). It encapsulates market routing, provider management, caching, and failover logic behind a simple, consistent API.

## Architecture

```text
                          ┌──────────────────────────┐
                          │      facade.Service       │
                          │                           │
  GetKline(req) ────────→ │  1. Parse symbol → Market │
  GetSpot(req)            │  2. Route to Manager      │
  GetInstruments(req)     │  3. Manager selects        │
  GetProfile(req)         │     best Provider          │
  GetFinancial(req)       │  4. Provider calls Adapter │
  GetAnnouncements(req)   │  5. Return domain response │
                          └──────────────────────────┘
                                      │
                    ┌─────────────────┼──────────────────┐
                    ▼                 ▼                   ▼
            Manager[CN]        Manager[US]         Manager[Crypto]
            ├─ eastmoney(100)  ├─ yahoo(100)       ├─ binance(100)
            ├─ sina(75)        └───────────         ├─ okx(75)
            ├─ tencent(50)                          └───────────
            ├─ tushare(25)
            └─ xueqiu(1)
```

## Quick Start

```go
svc := facade.NewService()
defer svc.Close()

// K 线数据
klineResp, err := svc.GetKline(ctx, kline.Request{
    Symbol:    "000001.SZ",
    Timeframe: kline.Timeframe1d,
    StartTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
    EndTime:   time.Date(2024, 1, 31, 0, 0, 0, 0, time.Local),
})

// 实时行情
spotResp, err := svc.GetSpot(ctx, spot.Request{
    Symbols: []string{"000001.SZ", "600519.SH"},
})

// 证券列表
instResp, err := svc.GetInstruments(ctx, instrument.Request{
    PageSize: 100,
})
```

---

## API Reference

### Service Methods

| Method | Description | Markets |
|--------|-------------|---------|
| `GetKline(ctx, req)` | 获取 K 线数据 | CN, US, HK, Crypto |
| `GetKlineWithTrace(ctx, req)` | 获取 K 线数据（含追踪信息） | CN, US, HK, Crypto |
| `GetSpot(ctx, req)` | 获取实时行情 | CN, US, HK, Crypto |
| `GetSpotWithTrace(ctx, req)` | 获取实时行情（含追踪信息） | CN, US, HK, Crypto |
| `GetInstruments(ctx, req)` | 获取证券列表 | CN, US, HK, Crypto |
| `GetProfile(ctx, req)` | 获取个股档案 | CN |
| `GetFinancial(ctx, req)` | 获取财务数据 | CN |
| `GetAnnouncements(ctx, req)` | 获取公告新闻 | CN |
| `GetStats()` | 返回统计信息 | - |
| `Close()` | 释放资源 | - |

### Service Options

```go
// 启用指标收集
svc := facade.NewService(
    facade.WithMetrics(myCollector),
)
```

---

## Market Routing

Service 通过 `domain.Symbol.Parse()` 自动识别市场：

| Symbol Format | Parsed Market | Routed To |
|--------------|---------------|-----------|
| `000001.SZ` | CN | eastmoney → sina → tencent → tushare → xueqiu |
| `600519.SH` | CN | eastmoney → sina → tencent → tushare → xueqiu |
| `AAPL.US` | US | yahoo |
| `00700.HK.HKEX` | HK | eastmoneyhk |
| `BTCUSDT` | Crypto | binance → okx |

### GetInstruments Market Detection

`GetInstruments` 支持额外的 Market 参数用于市场路由：

| Market Value | Routed To |
|-------------|-----------|
| `""` (default) | CN |
| `"US"`, `"NASDAQ"`, `"NYSE"` | US |
| `"HK"`, `"HKEX"` | HK |
| `"CRYPTO"`, `"BINANCE"` | Crypto |
| `"USDT"` (quote asset) | Crypto |

---

## Data Provider Priority

### A 股 (CN)

| Domain | Provider 1 (100) | Provider 2 (75) | Provider 3 (50) | Provider 4 (25) | Provider 5 (1) |
|--------|------------------|------------------|------------------|------------------|----------------|
| K线 | eastmoney | sina | tencent | tushare | xueqiu |
| 行情 | sina | tencent | eastmoney | xueqiu | - |
| 证券列表 | eastmoney | tushare | cninfo | sse/szse | bse |
| 个股档案 | eastmoney | tushare | - | - | - |
| 财务数据 | eastmoney | tushare | - | - | - |
| 公告新闻 | eastmoney | cninfo | - | - | - |

### 美股 (US)

| Domain | Provider 1 (100) |
|--------|------------------|
| K线 | yahoo |
| 行情 | yahoo |
| 证券列表 | yahoo |

### 港股 (HK)

| Domain | Provider 1 (100) |
|--------|------------------|
| K线 | eastmoneyhk |
| 行情 | eastmoneyhk |
| 证券列表 | eastmoneyhk |

### 加密货币 (Crypto)

| Domain | Provider 1 (100) | Provider 2 (75) |
|--------|------------------|------------------|
| K线 | binance | okx |
| 行情 | binance | okx |
| 证券列表 | binance | okx |

---

## Caching

Service 使用两级缓存策略：

| Data Type | L1 Cache (Hot) | L2 Cache (TTL) |
|-----------|----------------|----------------|
| K线 | 1 min | 5 min |
| 行情 | 1 min | 10 sec |
| 证券列表/档案/财务/公告 | 1 min | 1 hour |

---

## Testing Guidelines

### Integration Tests

Facade 测试是真正的集成测试，需调用外部 API：

1. **Use `checkFacadeError`**: 统一处理 API 错误，优雅跳过不可控的外部故障
2. **Log sample data**: 使用 `t.Logf` 打印结果摘要，证明集成正常
3. **Cover all markets**: 测试覆盖 CN、US、HK、Crypto 四个市场
4. **Test error cases**: 验证无效 symbol 返回合适的错误信息

### Error Handling Pattern

```go
func checkFacadeError(t *testing.T, err error) {
    t.Helper()
    if err == nil {
        return
    }
    msg := err.Error()
    skipPatterns := []string{
        "timeout", "connection refused", "all providers failed",
        "403", "401", "429", "rate limit", "geo-restrict",
    }
    for _, p := range skipPatterns {
        if strings.Contains(strings.ToLower(msg), strings.ToLower(p)) {
            t.Skipf("Skipping: external API issue: %v", err)
            return
        }
    }
    t.Fatalf("Unexpected error: %v", err)
}
```

### Symbol Format Notes

- **A 股**: `000001.SZ`, `600519.SH`
- **美股**: `AAPL.US`, `MSFT.US.NASDAQ`
- **港股**: `00700.HK.HKEX` (注意：`00700.HK` 会被错误路由到 US 市场)
- **加密**: `BTCUSDT`, `ETHUSDT`

---

## Adding New Market Support

1. Create adapter(s) in `adapters/<provider>/`
2. Register in `facade/service.go` `initManagers()` with appropriate priority
3. Add market routing logic in `getMarketFromSymbol()` if needed
4. Add `GetInstruments` market detection in the `if req.Market != ""` block if needed
5. Add integration tests for the new market
6. Update this README's Provider Priority tables

## Client Constructor Pattern

When creating clients in `initManagers()`, always use `WithHTTPClient` option:

```go
// Correct
eastmoneyclient.NewClient(eastmoneyclient.WithHTTPClient(s.httpClient))

// Wrong — will not compile (type mismatch)
eastmoneyclient.NewClient(s.httpClient)
```
