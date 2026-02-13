package main

import (
	"context"
	"fmt"
	"time"

	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/domain/spot"
	"github.com/souloss/quantds/facade"
	"github.com/souloss/quantds/manager"
)

func main() {
	// Initialize the service with Metrics enabled
	// We use a MemoryCollector to store metrics in memory.
	collector := manager.NewMemoryCollector()
	svc := facade.NewService(
		facade.WithMetrics(collector),
	)

	ctx := context.Background()

	// Example 1: Get Real-time Spot Quote
	fmt.Println("=== 1. Getting Spot Quote (First Call) ===")
	spotReq := spot.Request{
		Symbols: []string{"600000.SH", "000001.SZ"},
	}
	spotResp, err := svc.GetSpot(ctx, spotReq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		printSpot(spotResp)
	}

	// Print Stats - Should show 1 fetch, 0 cache hits
	printStats(svc.GetStats(), "After First Spot Call")

	// Example 2: Get Spot Quote Again (Demonstrate Caching)
	fmt.Println("\n=== 2. Getting Spot Quote (Second Call - Should hit Cache) ===")
	spotResp2, err := svc.GetSpot(ctx, spotReq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		printSpot(spotResp2)
	}

	// Print Stats - Should show 2 fetches, 1 cache hit
	// The 'Duration' in metrics also reflects the low-level HTTP tracking latency.
	printStats(svc.GetStats(), "After Second Spot Call")

	// Example 3: Get K-Line Data
	fmt.Println("\n=== 3. Getting K-Line Data ===")
	klineReq := kline.Request{
		Symbol:    "600000.SH",
		Timeframe: kline.Timeframe1d,
		StartTime: time.Now().AddDate(0, 0, -20),
		EndTime:   time.Now(),
		Adjust:    kline.AdjustNone,
	}
	klineResp, err := svc.GetKline(ctx, klineReq)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		printKline(klineResp)
	}

	printStats(svc.GetStats(), "Final Stats")

	// Example 4: Get K-Line Data with Trace (Show low-level HTTP details)
	fmt.Println("\n=== 4. Getting K-Line Data with Trace ===")
	// Using a new request to avoid cache (time range slightly different)
	klineReqTrace := kline.Request{
		Symbol:    "000001.SZ",
		Timeframe: kline.Timeframe1d,
		StartTime: time.Now().AddDate(0, 0, -10),
		EndTime:   time.Now(),
		Adjust:    kline.AdjustNone,
	}
	_, trace, err := svc.GetKlineWithTrace(ctx, klineReqTrace)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		printTrace(trace)
	}
}

func printSpot(resp spot.Response) {
	for _, q := range resp.Quotes {
		fmt.Printf("Stock: %s (%s) Price: %.2f Change: %.2f%%\n",
			q.Name, q.Symbol, q.Latest, q.ChangeRate)
	}
}

func printKline(resp kline.Response) {
	bars := resp.Bars
	if len(bars) > 5 {
		bars = bars[len(bars)-5:]
	}
	fmt.Printf("Symbol: %s (Total %d bars)\n", resp.Symbol, len(resp.Bars))
	for _, k := range bars {
		fmt.Printf("  %s | O: %.2f C: %.2f | Vol: %.0f\n",
			k.Timestamp.Format("2006-01-02"), k.Open, k.Close, k.Volume)
	}
}

func printStats(stats manager.Stats, title string) {
	fmt.Printf("\n--- Metrics [%s] ---\n", title)
	fmt.Printf("Total Fetches:   %d\n", stats.TotalFetches)
	fmt.Printf("Cache Hits:      %d\n", stats.CacheHits)
	fmt.Printf("Success Fetches: %d\n", stats.SuccessFetches)
	fmt.Printf("Failed Fetches:  %d\n", stats.FailedFetches)
	fmt.Printf("Avg Latency:     %v\n", stats.AvgLatency)
	if len(stats.ByProvider) > 0 {
		fmt.Println("By Provider:")
		for name, p := range stats.ByProvider {
			fmt.Printf("  - %-10s: Fetches=%d Success=%d Latency=%v\n",
				name, p.Fetches, p.Success, time.Duration(p.Duration)/time.Duration(p.Fetches))
		}
	}
	fmt.Println("---------------------------")
}

func printTrace(trace *manager.RequestTrace) {
	if trace == nil {
		fmt.Println("No trace information available.")
		return
	}
	fmt.Printf("\n--- HTTP Request Trace [FetchID: %s] ---\n", trace.FetchID)
	fmt.Printf("Total Requests: %d\n", trace.TotalRequests())
	fmt.Printf("Total Duration: %v\n", trace.TotalTime)

	for i, req := range trace.Requests {
		fmt.Printf("\n[Request #%d] (Attempt: %d)\n", i+1, req.Attempt)
		fmt.Printf("  URL:      %s %s\n", req.Request.Method, req.Request.URL)
		fmt.Printf("  Status:   %d\n", req.Response.StatusCode)
		fmt.Printf("  Duration: %v\n", req.Duration)
		if req.IsError() {
			fmt.Printf("  Error:    %v\n", req.Error)
		} else {
			body := string(req.Response.Body)
			fmt.Printf("  Body:     %s\n", body)
		}
	}
	fmt.Println("-------------------------------------------")
}
