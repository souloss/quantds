package okx

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/souloss/quantds/clients/okx"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

// Adapter name
const Name = "okx"

// Supported markets for OKX adapter
var supportedMarkets = []domain.Market{domain.MarketCrypto}

// KlineAdapter adapts OKX candlestick data to domain kline
type KlineAdapter struct {
	client *okx.Client
}

// NewKlineAdapter creates a new K-line adapter
func NewKlineAdapter(client *okx.Client) *KlineAdapter {
	return &KlineAdapter{client: client}
}

// Name returns the adapter name
func (a *KlineAdapter) Name() string {
	return Name
}

// SupportedMarkets returns supported markets
func (a *KlineAdapter) SupportedMarkets() []domain.Market {
	return supportedMarkets
}

// CanHandle checks if the adapter can handle the symbol
func (a *KlineAdapter) CanHandle(symbol string) bool {
	var sym domain.Symbol
	if err := sym.Parse(symbol); err != nil {
		// Try OKX format directly (e.g., "BTC-USDT")
		return strings.Contains(symbol, "-")
	}
	for _, m := range supportedMarkets {
		if sym.Market == m {
			return true
		}
	}
	return false
}

// Fetch retrieves K-line data from OKX
func (a *KlineAdapter) Fetch(ctx context.Context, _ request.Client, req kline.Request) (kline.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	// Convert symbol to OKX format (e.g., "BTCUSDT" → "BTC-USDT")
	instID := toOKXInstID(req.Symbol)

	params := &okx.CandlestickRequest{
		InstID: instID,
		Bar:    toOKXBar(req.Timeframe),
		Limit:  500,
	}

	candles, record, err := a.client.GetCandlesticks(ctx, params)
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(candles))
	for _, c := range candles {
		bar, ok := parseCandlestick(c)
		if !ok {
			continue
		}
		bars = append(bars, bar)
	}

	trace.Finish()
	return kline.Response{
		Symbol: req.Symbol,
		Bars:   bars,
		Source: Name,
	}, trace, nil
}

// parseCandlestick converts OKX candlestick data [ts, o, h, l, c, vol, volCcy, volCcyQuote, confirm]
func parseCandlestick(c okx.CandlestickResponse) (kline.Bar, bool) {
	if len(c) < 7 {
		return kline.Bar{}, false
	}

	ts, err := strconv.ParseInt(c[0], 10, 64)
	if err != nil {
		return kline.Bar{}, false
	}

	open, _ := strconv.ParseFloat(c[1], 64)
	high, _ := strconv.ParseFloat(c[2], 64)
	low, _ := strconv.ParseFloat(c[3], 64)
	close_, _ := strconv.ParseFloat(c[4], 64)
	vol, _ := strconv.ParseFloat(c[5], 64)
	turnover, _ := strconv.ParseFloat(c[6], 64)

	return kline.Bar{
		Timestamp: time.UnixMilli(ts),
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close_,
		Volume:    vol,
		Turnover:  turnover,
	}, true
}

// toOKXInstID converts a domain symbol to OKX instrument ID format
// e.g., "BTCUSDT" → "BTC-USDT", already "BTC-USDT" → "BTC-USDT"
func toOKXInstID(symbol string) string {
	// Already in OKX format
	if strings.Contains(symbol, "-") {
		return symbol
	}
	// Try common quote currencies
	quoteCurrencies := []string{"USDT", "USDC", "USD", "BTC", "ETH"}
	for _, q := range quoteCurrencies {
		if strings.HasSuffix(strings.ToUpper(symbol), q) {
			base := symbol[:len(symbol)-len(q)]
			if base != "" {
				return strings.ToUpper(base) + "-" + q
			}
		}
	}
	return symbol
}

// fromOKXInstID converts an OKX instrument ID to domain format
// e.g., "BTC-USDT" → "BTCUSDT"
func fromOKXInstID(instID string) string {
	return strings.ReplaceAll(instID, "-", "")
}

// toOKXBar converts domain timeframe to OKX bar format
func toOKXBar(tf kline.Timeframe) string {
	switch tf {
	case kline.Timeframe1m:
		return "1m"
	case kline.Timeframe5m:
		return "5m"
	case kline.Timeframe15m:
		return "15m"
	case kline.Timeframe30m:
		return "30m"
	case kline.Timeframe60m:
		return "1H"
	case kline.Timeframe1d:
		return "1D"
	case kline.Timeframe1w:
		return "1W"
	case kline.Timeframe1M:
		return "1M"
	default:
		return "1D"
	}
}

var _ manager.Provider[kline.Request, kline.Response] = (*KlineAdapter)(nil)
