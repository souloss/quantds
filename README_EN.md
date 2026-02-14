# quantds - Multi-Market Multi-Asset Data Service

[简体中文](README.md) | English

`quantds` is a unified financial data service module that provides standardized access to data across multiple markets (A-shares, HK stocks, US stocks, etc.) and multiple assets (Stocks, Funds, Bonds, etc.). It encapsulates multiple data sources (EastMoney, Sina, Tencent, Xueqiu, Tushare, etc.) and provides automatic failover, load balancing, and caching capabilities.

## Features

*   **Unified Interface**: Uses a unified Domain Model to access data from different data sources.
*   **Multi-Source Support**: Built-in clients and adapters for mainstream financial data sources.
*   **High Availability**: Supports automatic failover for multiple providers. When one data source is unavailable, it automatically tries the next one.
*   **High Performance**: Built-in two-level caching (memory + distributed cache interface) to reduce network requests and improve response speed.
*   **Easy to Extend**: Adopts Facade - Manager - Adapter - Client layered architecture, making it easy to add new data sources.

## Supported Data Sources

<!-- START_STATUS_BADGES -->

![](https://img.shields.io/badge/Binance-%E2%9C%93%20Crypto-brightgreen) ![](https://img.shields.io/badge/BSE-%E2%9C%93%20A%E8%82%A1-brightgreen) ![](https://img.shields.io/badge/Cninfo-%E2%9C%93%20A%E8%82%A1-brightgreen) ![](https://img.shields.io/badge/EastMoney-%E2%9C%93%20A%E8%82%A1-brightgreen) ![](https://img.shields.io/badge/EastMoneyHK-%E2%9C%93%20%E6%B8%AF%E8%82%A1-brightgreen) ![](https://img.shields.io/badge/Okx-%E2%9C%93%20Ready-brightgreen) ![](https://img.shields.io/badge/Sina-%E2%9C%93%20A%E8%82%A1%20%7C%20%E6%B8%AF%E8%82%A1-brightgreen) ![](https://img.shields.io/badge/SSE-%E2%9C%93%20A%E8%82%A1-brightgreen) ![](https://img.shields.io/badge/SZSE-%E2%9C%93%20A%E8%82%A1-brightgreen) ![](https://img.shields.io/badge/Tencent-%E2%9C%93%20A%E8%82%A1%20%7C%20%E6%B8%AF%E8%82%A1%20%7C%20%E7%BE%8E%E8%82%A1-brightgreen) ![](https://img.shields.io/badge/Tushare-%E2%9C%93%20A%E8%82%A1-brightgreen) ![](https://img.shields.io/badge/Xueqiu-%F0%9F%9F%A1%20Beta-yellow) ![](https://img.shields.io/badge/Yahoo-%E2%9C%93%20A%E8%82%A1%20%7C%20%E6%B8%AF%E8%82%A1%20%7C%20%E7%BE%8E%E8%82%A1-brightgreen) 

<!-- END_STATUS_BADGES -->

<!-- START_SUPPORTED_TABLE -->

| 数据源 (Provider) | K线 (Kline) | 实时行情 (Spot) | 证券列表 (Instrument) | 证券详情 (Profile) | 财务数据 (Financial) | 公告资讯 (News) |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: |
| **Binance** | ✅ | ✅ | ✅ | - | - | - |
| **BSE** | - | - | ✅ | - | - | - |
| **Cninfo** | - | - | ✅ | - | - | ✅ |
| **EastMoney** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **EastMoneyHK** | ✅ | ✅ | ✅ | - | - | - |
| **Okx** | ✅ | ✅ | ✅ | - | - | - |
| **Sina** | ✅ | ✅ | - | - | - | - |
| **SSE** | - | - | ✅ | - | - | - |
| **SZSE** | - | - | ✅ | - | - | - |
| **Tencent** | ✅ | ✅ | - | - | - | - |
| **Tushare** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Xueqiu** | ✅ | ✅ | ✅ | ✅ | - | - |
| **Yahoo** | ✅ | ✅ | ✅ | - | - | - |


<!-- END_SUPPORTED_TABLE -->

## Directory Structure

```
quantds/
├── adapters/          # Adapter Layer: Converts data from various sources into unified domain models
│   ├── eastmoney/     # EastMoney Adapter
│   ├── sina/          # Sina Adapter
│   └── ...
├── clients/           # Client Layer: Direct interface with third-party APIs
│   ├── eastmoney/
│   ├── sina/
│   └── ...
├── domain/            # Domain Layer: Defines common data models and interfaces
│   ├── kline/         # Kline Data Model (Bar, Request, Response)
│   ├── spot/          # Spot Quote Model (Quote, Request, Response)
│   ├── instrument/    # Instrument List Model
│   ├── profile/       # Stock Profile Model
│   ├── financial/     # Financial Data Model
│   └── announcement/  # Announcement Model
├── facade/            # Facade Layer: Unified external entry point (Service)
├── manager/           # Manager Layer: Responsible for Provider management, routing, caching, monitoring
├── request/           # Basic HTTP Client Encapsulation
└── example/           # Usage Examples
```

## Usage Examples

### 1. Installation

```bash
go get github.com/souloss/quantds
```

### 2. Quick Start

The following code demonstrates how to initialize the service and retrieve real-time quotes and K-line data. For the complete code, please refer to `example/main.go`.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/facade"
)

func main() {
	// 1. Initialize Service
	// The Service automatically loads all registered data source adapters and configures default priority and caching strategies.
	svc := facade.NewService()
	ctx := context.Background()

	// 2. Get Real-time Spot Quote
	fmt.Println("=== Get Real-time Spot Quote ===")
	spotReq := spot.Request{
		Symbols: []string{"600000.SH", "000001.SZ"}, // Supports multiple formats, e.g., sh600000, 600000.SH
	}
	spotResp, err := svc.GetSpot(ctx, spotReq)
	if err != nil {
		panic(err)
	}

	for _, q := range spotResp.Quotes {
		fmt.Printf("Stock: %s (%s) Price: %.2f Change: %.2f%%\n",
			q.Name, q.Symbol, q.Latest, q.ChangeRate)
	}

	// 3. Get K-Line Data (Candlestick)
	fmt.Println("\n=== Get Daily K-Line Data ===")
	klineReq := kline.Request{
		Symbol:    "600000.SH",
		Timeframe: kline.Timeframe1d,           // Period: Daily
		StartTime: time.Now().AddDate(0, 0, -10), // Last 10 days
		EndTime:   time.Now(),
		Adjust:    kline.AdjustNone,            // Adjustment: None
	}
	klineResp, err := svc.GetKline(ctx, klineReq)
	if err != nil {
		panic(err)
	}

	for _, bar := range klineResp.Bars {
		fmt.Printf("Date: %s Open: %.2f Close: %.2f Vol: %.0f\n",
			bar.Timestamp.Format("2006-01-02"), bar.Open, bar.Close, bar.Volume)
	}
}
```

### 3. Advanced Features: Metrics and Caching

`quantds` has built-in Metrics Collector and multi-level caching mechanisms.

```go
import (
	"github.com/souloss/quantds/facade"
	"github.com/souloss/quantds/manager"
)

func main() {
	// Enable Memory Metrics Collector
	collector := manager.NewMemoryCollector()
	svc := facade.NewService(
		facade.WithMetrics(collector),
	)

	// ... Execute data requests ...

	// Get and print statistics
	stats := svc.GetStats()
	fmt.Printf("Total Fetches: %d\n", stats.TotalFetches)
	fmt.Printf("Cache Hits:    %d\n", stats.CacheHits)
	fmt.Printf("Avg Latency:   %v\n", stats.AvgLatency)

	// Get request trace details
	_, trace, _ := svc.GetKlineWithTrace(context.Background(), req)
	if trace != nil {
		fmt.Printf("Request ID: %s, Total Duration: %v\n", trace.FetchID, trace.TotalTime)
		for _, r := range trace.Requests {
			fmt.Printf("HTTP %s %s -> Status %d\n", r.Request.Method, r.Request.URL, r.Response.StatusCode)
		}
	}
}
```

## Architecture

`quantds` adopts a layered architecture design:

1.  **Facade Layer**: `facade.Service` is the only external entry point, shielding complex internal scheduling logic. Users only need to call simple interfaces like `GetSpot`, `GetKline`.
2.  **Manager Layer**: The core scheduling center.
    *   **Router/Selector**: Selects appropriate data source Providers based on priority or weight.
    *   **Failover**: Automatically degrades to the next available Provider when a high-priority Provider fails.
    *   **Cache**: Handles caching logic uniformly, supporting request-level caching.
3.  **Adapter Layer**: Implements the unified `Provider` interface, converting heterogeneous data from different sources into standard Domain Models.
4.  **Client Layer**: Pure HTTP API encapsulation, responsible for handling request signatures, protocol parsing, etc., containing no business logic.

## Naming Conventions (Current)

The project currently uses the following domain naming conventions:

*   **kline**: K-line data / Candlestick (corresponds to `domain/kline`)
*   **spot**: Real-time quote / Snapshot (corresponds to `domain/spot`)
*   **instrument**: Instrument list / Basic info (corresponds to `domain/instrument`)
*   **profile**: Stock profile / Depth info (corresponds to `domain/profile`)
*   **financial**: Financial statements data (corresponds to `domain/financial`)
*   **announcement**: Announcements and News (corresponds to `domain/announcement`)
