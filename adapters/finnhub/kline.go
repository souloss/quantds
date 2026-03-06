package finnhub

import (
	"context"
	"time"

	"github.com/souloss/quantds/clients/finnhub"
	"github.com/souloss/quantds/domain"
	"github.com/souloss/quantds/domain/kline"
	"github.com/souloss/quantds/manager"
	"github.com/souloss/quantds/request"
)

const Name = "finnhub"

var supportedMarkets = []domain.Market{domain.MarketUS, domain.MarketForex, domain.MarketCrypto}

type KlineAdapter struct {
	client *finnhub.Client
}

func NewKlineAdapter(client *finnhub.Client) *KlineAdapter {
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

func (a *KlineAdapter) Fetch(ctx context.Context, _ request.Client, req kline.Request) (kline.Response, *manager.RequestTrace, error) {
	trace := manager.NewRequestTrace(Name)

	var sym domain.Symbol
	if err := sym.Parse(req.Symbol); err != nil {
		return kline.Response{}, trace, err
	}

	from := req.StartTime
	to := req.EndTime
	if from.IsZero() {
		from = time.Now().AddDate(-1, 0, 0)
	}
	if to.IsZero() {
		to = time.Now()
	}

	params := &finnhub.CandleParams{
		Symbol:     sym.Code,
		Resolution: finnhub.ToResolution(string(req.Timeframe)),
		From:       from.Unix(),
		To:         to.Unix(),
	}

	var result *finnhub.CandleResult
	var record *request.Record
	var err error

	switch sym.Market {
	case domain.MarketForex:
		params.Symbol = "OANDA:" + sym.Code
		result, record, err = a.client.GetForexCandles(ctx, params)
	case domain.MarketCrypto:
		params.Symbol = "BINANCE:" + sym.Code
		result, record, err = a.client.GetCryptoCandles(ctx, params)
	default:
		result, record, err = a.client.GetStockCandles(ctx, params)
	}
	trace.AddRequest(record)

	if err != nil {
		return kline.Response{}, trace, err
	}

	bars := make([]kline.Bar, 0, len(result.Candles))
	for _, c := range result.Candles {
		bars = append(bars, kline.Bar{
			Timestamp: time.Unix(c.Timestamp, 0),
			Open:      c.Open,
			High:      c.High,
			Low:       c.Low,
			Close:     c.Close,
			Volume:    c.Volume,
		})
	}

	trace.Finish()
	return kline.Response{
		Symbol: req.Symbol,
		Bars:   bars,
		Source: Name,
	}, trace, nil
}

var _ manager.Provider[kline.Request, kline.Response] = (*KlineAdapter)(nil)
