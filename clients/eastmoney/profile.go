package eastmoney

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/souloss/quantds/request"
)

const ProfileAPI = "https://push2.eastmoney.com/api/qt/stock/get"
const ProfileFields = "f57,f58,f116,f117,f84,f85,f162,f167,f55,f92,f127,f59,f26,f189,f43,f46,f44,f45,f60,f47,f48,f168,f170,f163,f164,f169,f191,f171"

// ProfileParams represents parameters for security profile request
type ProfileParams struct {
	Symbol string
}

// ProfileResult represents the security profile result
type ProfileResult struct {
	Code           string
	Name           string
	Currency       string
	ListingDate    string
	Industry       string
	LatestPrice    float64
	Open           float64
	High           float64
	Low            float64
	PreClose       float64
	Volume         float64
	Amount         float64
	TurnoverRate   float64
	ChangePct      float64
	VolumeRatio    float64
	Amplitude      float64
	TotalMarketCap float64
	FloatMarketCap float64
	TotalShares    float64
	FloatShares    float64
	PEDynamic      float64
	PEStatic       float64
	PETTM          float64
	PBRatio        float64
	PSTTM          float64
	EPS            float64
	NAVPS          float64
}

// DetailParams is an alias for ProfileParams
type DetailParams = ProfileParams

// GetProfile retrieves detailed information about a security
func (c *Client) GetProfile(ctx context.Context, params *ProfileParams) (*ProfileResult, *request.Record, error) {
	if params.Symbol == "" {
		return nil, nil, fmt.Errorf("symbol required")
	}

	secid, err := toEastMoneySecid(params.Symbol)
	if err != nil {
		return nil, nil, err
	}

	url := fmt.Sprintf("%s?secid=%s&ut=fa5fd1943c7b386f172d6893dbfba10b&fltt=2&invt=2&fields=%s",
		PushURL, secid, ProfileFields)

	req := request.Request{
		Method: "GET",
		URL:    url,
		Headers: map[string]string{
			"Host":       "push2.eastmoney.com",
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Referer":    "https://quote.eastmoney.com/",
		},
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseProfileResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

// GetStockDetail is an alias for GetProfile
func (c *Client) GetStockDetail(ctx context.Context, params *DetailParams) (*ProfileResult, *request.Record, error) {
	return c.GetProfile(ctx, params)
}

type profileResponse struct {
	Data map[string]interface{} `json:"data"`
}

func parseProfileResponse(body []byte) (*ProfileResult, error) {
	var resp profileResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if resp.Data == nil {
		return nil, fmt.Errorf("empty data response")
	}

	return &ProfileResult{
		Code:           getDetailStr(resp.Data, "f57"),
		Name:           getDetailStr(resp.Data, "f58"),
		Currency:       getDetailStr(resp.Data, "f59"),
		ListingDate:    getDetailStr(resp.Data, "f26"),
		Industry:       getDetailStr(resp.Data, "f127"),
		LatestPrice:    getDetailFlt(resp.Data, "f43"),
		Open:           getDetailFlt(resp.Data, "f46"),
		High:           getDetailFlt(resp.Data, "f44"),
		Low:            getDetailFlt(resp.Data, "f45"),
		PreClose:       getDetailFlt(resp.Data, "f60"),
		Volume:         getDetailFlt(resp.Data, "f47"),
		Amount:         getDetailFlt(resp.Data, "f48"),
		TurnoverRate:   getDetailFlt(resp.Data, "f168"),
		ChangePct:      getDetailFlt(resp.Data, "f170"),
		VolumeRatio:    getDetailFlt(resp.Data, "f191"),
		Amplitude:      getDetailFlt(resp.Data, "f171"),
		TotalMarketCap: getDetailFlt(resp.Data, "f116"),
		FloatMarketCap: getDetailFlt(resp.Data, "f117"),
		TotalShares:    getDetailFlt(resp.Data, "f84"),
		FloatShares:    getDetailFlt(resp.Data, "f85"),
		PEDynamic:      getDetailFlt(resp.Data, "f162"),
		PEStatic:       getDetailFlt(resp.Data, "f163"),
		PETTM:          getDetailFlt(resp.Data, "f164"),
		PBRatio:        getDetailFlt(resp.Data, "f167"),
		PSTTM:          getDetailFlt(resp.Data, "f169"),
		EPS:            getDetailFlt(resp.Data, "f55"),
		NAVPS:          getDetailFlt(resp.Data, "f92"),
	}, nil
}

func getDetailStr(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		switch val := v.(type) {
		case string:
			return val
		case float64:
			return fmt.Sprintf("%.0f", val)
		}
	}
	return ""
}

func getDetailFlt(data map[string]interface{}, key string) float64 {
	if v, ok := data[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case string:
			var f float64
			fmt.Sscanf(val, "%f", &f)
			return f
		}
	}
	return 0
}
