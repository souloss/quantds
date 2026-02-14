# quantds - 多市场多资产数据服务

简体中文 | [English](README_EN.md)

`quantds` 是一个统一的金融数据服务模块，提供对多市场（A股、港股、美股等）和多资产（股票、基金、债券等）数据的标准化访问。它封装了多个数据源（EastMoney, Sina, Tencent, Xueqiu, Tushare 等），提供自动故障转移（Failover）、负载均衡和缓存功能。

## 功能特性

*   **统一接口**: 使用统一的领域模型（Domain Model）访问不同数据源的数据。
*   **多数据源支持**: 内置主流财经数据源的客户端和适配器。
*   **高可用性**: 支持多 Provider 自动切换（Failover），当某个数据源不可用时自动尝试下一个。
*   **高性能**: 内置两级缓存（内存 + 分布式缓存接口），减少网络请求，提高响应速度。
*   **易扩展**: 采用 Facade - Manager - Adapter - Client 分层架构，易于添加新的数据源。

## 支持的数据源

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

## 目录结构

```
quantds/
├── adapters/          # 适配器层：将各数据源的数据转换为统一领域模型
│   ├── eastmoney/     # 东方财富适配器
│   ├── sina/          # 新浪财经适配器
│   └── ...
├── clients/           # 客户端层：直接对接第三方 API
│   ├── eastmoney/
│   ├── sina/
│   └── ...
├── domain/            # 领域层：定义通用数据模型和接口
│   ├── kline/         # K线数据模型 (Bar, Request, Response)
│   ├── spot/          # 实时行情模型 (Quote, Request, Response)
│   ├── instrument/    # 证券列表模型
│   ├── profile/       # 证券详情模型
│   ├── financial/     # 财务数据模型
│   └── announcement/  # 公告资讯模型
├── facade/            # 外观层：对外统一入口 (Service)
├── manager/           # 管理层：负责 Provider 管理、路由、缓存、监控
├── request/           # 基础 HTTP 客户端封装
└── example/           # 使用示例
```

## 使用示例

### 1. 安装

```bash
go get github.com/souloss/quantds
```

### 2. 快速开始

以下代码展示了如何初始化服务并获取实时行情和 K 线数据。完整代码请参考 `example/main.go`。

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
	// 1. 初始化服务
	// Service 会自动加载所有已注册的数据源适配器，并配置默认的优先级和缓存策略
	svc := facade.NewService()
	ctx := context.Background()

	// 2. 获取实时行情 (Spot Quote)
	fmt.Println("=== 获取实时行情 ===")
	spotReq := spot.Request{
		Symbols: []string{"600000.SH", "000001.SZ"}, // 支持多种格式，如 sh600000, 600000.SH
	}
	spotResp, err := svc.GetSpot(ctx, spotReq)
	if err != nil {
		panic(err)
	}

	for _, q := range spotResp.Quotes {
		fmt.Printf("股票: %s (%s) 价格: %.2f 涨跌幅: %.2f%%\n",
			q.Name, q.Symbol, q.Latest, q.ChangeRate)
	}

	// 3. 获取 K 线数据 (Kline/Candlestick)
	fmt.Println("\n=== 获取日 K 线数据 ===")
	klineReq := kline.Request{
		Symbol:    "600000.SH",
		Timeframe: kline.Timeframe1d,           // 周期：日线
		StartTime: time.Now().AddDate(0, 0, -10), // 最近10天
		EndTime:   time.Now(),
		Adjust:    kline.AdjustNone,            // 复权：不复权
	}
	klineResp, err := svc.GetKline(ctx, klineReq)
	if err != nil {
		panic(err)
	}

	for _, bar := range klineResp.Bars {
		fmt.Printf("日期: %s 开: %.2f 收: %.2f 量: %.0f\n",
			bar.Timestamp.Format("2006-01-02"), bar.Open, bar.Close, bar.Volume)
	}
}
```

### 3. 高级功能：监控指标与缓存

`quantds` 内置了指标收集器（Metrics Collector）和多级缓存机制。

```go
import (
	"github.com/souloss/quantds/facade"
	"github.com/souloss/quantds/manager"
)

func main() {
	// 启用内存指标收集器
	collector := manager.NewMemoryCollector()
	svc := facade.NewService(
		facade.WithMetrics(collector),
	)

	// ... 执行数据请求 ...

	// 获取并打印统计信息
	stats := svc.GetStats()
	fmt.Printf("Total Fetches: %d\n", stats.TotalFetches)
	fmt.Printf("Cache Hits:    %d\n", stats.CacheHits)
	fmt.Printf("Avg Latency:   %v\n", stats.AvgLatency)

	// 获取请求追踪详情
	_, trace, _ := svc.GetKlineWithTrace(context.Background(), req)
	if trace != nil {
		fmt.Printf("Request ID: %s, Total Duration: %v\n", trace.FetchID, trace.TotalTime)
		for _, r := range trace.Requests {
			fmt.Printf("HTTP %s %s -> Status %d\n", r.Request.Method, r.Request.URL, r.Response.StatusCode)
		}
	}
}
```

## 架构说明

`quantds` 采用分层架构设计：

1.  **Facade (外观层)**: `facade.Service` 是对外的唯一入口，屏蔽了内部复杂的调度逻辑。用户只需调用 `GetSpot`, `GetKline` 等简单接口。
2.  **Manager (管理层)**: 核心调度中心。
    *   **Router/Selector**: 根据优先级（Priority）或权重选择合适的数据源 Provider。
    *   **Failover**: 当高优先级 Provider 失败时，自动降级到下一个可用 Provider。
    *   **Cache**: 统一处理缓存逻辑，支持请求级缓存。
3.  **Adapter (适配器层)**: 实现统一的 `Provider` 接口，将不同数据源的异构数据转换为标准 Domain Model。
4.  **Client (客户端层)**: 纯粹的 HTTP API 封装，负责处理请求签名、协议解析等，不包含业务逻辑。

## 命名规范 (Current)

目前项目主要使用以下领域命名：

*   **kline**: K线数据 / 蜡烛图 (对应 `domain/kline`)
*   **spot**: 实时行情 / 快照 (对应 `domain/spot`)
*   **instrument**: 证券列表 / 基础信息 (对应 `domain/instrument`)
*   **profile**: 证券深度资料 (对应 `domain/profile`)
*   **financial**: 财务报表数据 (对应 `domain/financial`)
*   **announcement**: 公告与新闻 (对应 `domain/announcement`)
