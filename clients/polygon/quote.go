package polygon

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

type SnapshotParams struct {
	Tickers []string
}

type SnapshotResult struct {
	Tickers []SnapshotData
	Count   int
}

type SnapshotData struct {
	Ticker string
	Day    SnapshotBar
	Min    SnapshotBar
	PrevDay SnapshotBar
	Change  float64
	ChangePercent float64
	Updated int64
}

type SnapshotBar struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
	VWAP   float64
}

func (c *Client) GetSnapshot(ctx context.Context, params *SnapshotParams) (*SnapshotResult, *request.Record, error) {
	url := fmt.Sprintf("%s%s?apiKey=%s", BaseURL, SnapshotAPI, c.apiKey)

	if len(params.Tickers) > 0 {
		tickers := ""
		for i, t := range params.Tickers {
			if i > 0 {
				tickers += ","
			}
			tickers += t
		}
		url += fmt.Sprintf("&tickers=%s", tickers)
	}

	req := request.Request{
		Method:  "GET",
		URL:     url,
		Headers: DefaultHeaders,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseSnapshotResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

type polygonSnapshotResponse struct {
	Status  string `json:"status"`
	Tickers []struct {
		Ticker string `json:"ticker"`
		Day    struct {
			O  float64 `json:"o"`
			H  float64 `json:"h"`
			L  float64 `json:"l"`
			C  float64 `json:"c"`
			V  float64 `json:"v"`
			VW float64 `json:"vw"`
		} `json:"day"`
		Min struct {
			O  float64 `json:"o"`
			H  float64 `json:"h"`
			L  float64 `json:"l"`
			C  float64 `json:"c"`
			V  float64 `json:"v"`
			VW float64 `json:"vw"`
		} `json:"min"`
		PrevDay struct {
			O  float64 `json:"o"`
			H  float64 `json:"h"`
			L  float64 `json:"l"`
			C  float64 `json:"c"`
			V  float64 `json:"v"`
			VW float64 `json:"vw"`
		} `json:"prevDay"`
		TodaysChange    float64 `json:"todaysChange"`
		TodaysChangePer float64 `json:"todaysChangePerc"`
		Updated         int64   `json:"updated"`
	} `json:"tickers"`
}

func parseSnapshotResponse(body []byte) (*SnapshotResult, error) {
	var resp polygonSnapshotResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	tickers := make([]SnapshotData, 0, len(resp.Tickers))
	for _, t := range resp.Tickers {
		tickers = append(tickers, SnapshotData{
			Ticker:        t.Ticker,
			Day:           SnapshotBar{Open: t.Day.O, High: t.Day.H, Low: t.Day.L, Close: t.Day.C, Volume: t.Day.V, VWAP: t.Day.VW},
			Min:           SnapshotBar{Open: t.Min.O, High: t.Min.H, Low: t.Min.L, Close: t.Min.C, Volume: t.Min.V, VWAP: t.Min.VW},
			PrevDay:       SnapshotBar{Open: t.PrevDay.O, High: t.PrevDay.H, Low: t.PrevDay.L, Close: t.PrevDay.C, Volume: t.PrevDay.V, VWAP: t.PrevDay.VW},
			Change:        t.TodaysChange,
			ChangePercent: t.TodaysChangePer,
			Updated:       t.Updated,
		})
	}

	return &SnapshotResult{
		Tickers: tickers,
		Count:   len(tickers),
	}, nil
}
