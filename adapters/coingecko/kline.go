package coingecko

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/souloss/quantds/clients/coingecko"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
)

// KlineAdapter adapts CoinGecko market/chart data to kline domain
type KlineAdapter struct {
	client *coingecko.Client
}

// NewKlineAdapter creates a new kline adapter
func NewKlineAdapter(client *coingecko.Client) *KlineAdapter {
	return &KlineAdapter{client: client}
}

func (a *KlineAdapter) Name() string                      { return Name }
func (a *KlineAdapter) SupportedMarkets() []domain.Market { return supportedMarkets }

func (a *KlineAdapter) CanHandle(symbol string) bool {
	var sym domain.Symbol
	if err := sym.Parse(symbol); err != nil {
		return false
	}
	for _, m := range supportedMarkets {
		if sym.Market == m {
			return true
		}
	}
	return false
}

func (a *KlineAdapter) Fetch(ctx context.Context, req kline.Request) (kline.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	var sym domain.Symbol
	if err := sym.Parse(req.Symbol); err != nil {
		return kline.Response{}, trace, fmt.Errorf("invalid symbol: %s", req.Symbol)
	}

	coinID := symbolToCoinID(sym.Code)
	if coinID == "" {
		return kline.Response{}, trace, fmt.Errorf("unsupported coin: %s", sym.Code)
	}

	quoteCurrency := "usd"
	if sym.Exchange != "" {
		quoteCurrency = strings.ToLower(string(sym.Exchange))
	}

	// Calculate days from time range
	days := timeframeToDays(req.Timeframe, req.StartTime, req.EndTime)

	result, record, err := a.client.GetMarketChart(ctx, &coingecko.MarketChartRequest{
		ID:         coinID,
		VsCurrency: quoteCurrency,
		Days:       days,
		Interval:   timeframeToInterval(req.Timeframe),
	})
	trace.AddRequest(record)
	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Prices))
	for i, price := range result.Prices {
		if i >= len(result.TotalVolumes) {
			break
		}
		ts := time.UnixMilli(int64(price[0]))
		bars = append(bars, kline.Bar{
			Timestamp: ts,
			Open:      price[1], // CoinGecko simple mode only provides close prices
			High:      price[1],
			Low:       price[1],
			Close:     price[1],
			Volume:    result.TotalVolumes[i][1],
		})
	}

	trace.Finish()
	return kline.Response{
		Symbol:      req.Symbol,
		Bars:        bars,
		Source:      Name,
		DataVersion: 1,
	}, trace, nil
}

func timeframeToDays(tf kline.Timeframe, start, end time.Time) string {
	if !start.IsZero() && !end.IsZero() {
		days := int(end.Sub(start).Hours() / 24)
		if days <= 1 {
			return "1"
		}
		return fmt.Sprintf("%d", days)
	}
	switch tf {
	case kline.Timeframe1d:
		return "1"
	case kline.Timeframe1w, kline.Timeframe1M:
		return "30"
	default:
		return "30"
	}
}

func timeframeToInterval(tf kline.Timeframe) string {
	switch tf {
	case kline.Timeframe1d:
		return ""
	case kline.Timeframe1w, kline.Timeframe1M:
		return "daily"
	default:
		return "daily"
	}
}

var _ manager.Provider[kline.Request, kline.Response] = (*KlineAdapter)(nil)
