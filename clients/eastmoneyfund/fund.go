package eastmoneyfund

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/souloss/quantds/request"
)

type FundListResult struct {
	Funds []FundInfo
	Count int
}

type FundInfo struct {
	Code     string
	Abbr     string
	Name     string
	Type     string
	FullName string
}

func (c *Client) GetFundList(ctx context.Context) (*FundListResult, *request.Record, error) {
	req := request.Request{
		Method:  "GET",
		URL:     FundListURL,
		Headers: DefaultHeaders,
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	result, err := parseFundListResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

func parseFundListResponse(body []byte) (*FundListResult, error) {
	content := string(body)
	// The response is: var r = [["000001","HXCZ","华夏成长","混合型-偏股","华夏成长混合"],...]
	re := regexp.MustCompile(`var r = (\[.*\])`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse fund list response")
	}

	var rawData [][]string
	if err := json.Unmarshal([]byte(matches[1]), &rawData); err != nil {
		return nil, err
	}

	funds := make([]FundInfo, 0, len(rawData))
	for _, item := range rawData {
		if len(item) < 5 {
			continue
		}
		funds = append(funds, FundInfo{
			Code:     item[0],
			Abbr:     item[1],
			Name:     item[2],
			Type:     item[3],
			FullName: item[4],
		})
	}

	return &FundListResult{
		Funds: funds,
		Count: len(funds),
	}, nil
}

type FundEstimateParams struct {
	Code string
}

type FundEstimateResult struct {
	Code      string
	Name      string
	EstNAV    string
	EstChange string
	EstTime   string
}

func (c *Client) GetFundEstimate(ctx context.Context, params *FundEstimateParams) (*FundEstimateResult, *request.Record, error) {
	url := fmt.Sprintf("%s/%s.js", FundEstimateURL, params.Code)

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

	result, err := parseFundEstimateResponse(resp.Body)
	if err != nil {
		return nil, record, err
	}

	return result, record, nil
}

func parseFundEstimateResponse(body []byte) (*FundEstimateResult, error) {
	content := string(body)
	// Response format: jsonpgz({"fundcode":"000001","name":"华夏成长混合","jzrq":"2024-01-01","dwjz":"1.234","gsz":"1.235","gszzl":"0.08","gztime":"2024-01-02 15:00"});
	content = strings.TrimPrefix(content, "jsonpgz(")
	content = strings.TrimSuffix(content, ");")

	var data map[string]string
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return nil, err
	}

	return &FundEstimateResult{
		Code:      data["fundcode"],
		Name:      data["name"],
		EstNAV:    data["gsz"],
		EstChange: data["gszzl"],
		EstTime:   data["gztime"],
	}, nil
}
