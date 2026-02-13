package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/souloss/quantds/request"
)

// API Endpoint Constants
const (
	EndpointCoinData = "/coins/%s"
	
	// Query Parameters
	ParamLocalization   = "localization"
	ParamTickers        = "tickers"
	ParamMarketData     = "market_data"
	ParamCommunityData  = "community_data"
	ParamDeveloperData  = "developer_data"
	ParamSparkline      = "sparkline"
)

// CoinDataRequest represents parameters for coin data request
type CoinDataRequest struct {
	ID            string // Coin ID (e.g. "bitcoin")
	Localization  bool   // Include all localized languages in response
	Tickers       bool   // Include tickers data
	MarketData    bool   // Include market data
	CommunityData bool   // Include community data
	DeveloperData bool   // Include developer data
	Sparkline     bool   // Include sparkline 7 days data
}

// CoinDataResponse represents the response structure
type CoinDataResponse struct {
	ID                 string                 `json:"id"`
	Symbol             string                 `json:"symbol"`
	Name               string                 `json:"name"`
	WebSlug            string                 `json:"web_slug"`
	AssetPlatformID    string                 `json:"asset_platform_id"`
	Platforms          map[string]string      `json:"platforms"`
	DetailPlatforms    map[string]interface{} `json:"detail_platforms"`
	BlockTimeInMinutes int32                  `json:"block_time_in_minutes"`
	HashingAlgorithm   string                 `json:"hashing_algorithm"`
	Categories         []string               `json:"categories"`
	Description        map[string]string      `json:"description"`
	Links              map[string]interface{} `json:"links"`
	Image              map[string]string      `json:"image"`
	CountryOrigin      string                 `json:"country_origin"`
	GenesisDate        string                 `json:"genesis_date"`
	SentimentVotesUp   float64                `json:"sentiment_votes_up_percentage"`
	SentimentVotesDown float64                `json:"sentiment_votes_down_percentage"`
	MarketCapRank      int32                  `json:"market_cap_rank"`
	CoingeckoRank      int32                  `json:"coingecko_rank"`
	CoingeckoScore     float64                `json:"coingecko_score"`
	DeveloperScore     float64                `json:"developer_score"`
	CommunityScore     float64                `json:"community_score"`
	LiquidityScore     float64                `json:"liquidity_score"`
	PublicInterestScore float64               `json:"public_interest_score"`
	MarketData         *MarketData            `json:"market_data"`
	CommunityData      *CommunityData         `json:"community_data"`
	DeveloperData      *DeveloperData         `json:"developer_data"`
	LastUpdated        string                 `json:"last_updated"`
}

type MarketData struct {
	CurrentPrice                           map[string]float64 `json:"current_price"`
	TotalValueLocked                       interface{}        `json:"total_value_locked"`
	McapToTvlRatio                         interface{}        `json:"mcap_to_tvl_ratio"`
	FDVToTvlRatio                          interface{}        `json:"fdv_to_tvl_ratio"`
	Roi                                    interface{}        `json:"roi"`
	Ath                                    map[string]float64 `json:"ath"`
	AthChangePercentage                    map[string]float64 `json:"ath_change_percentage"`
	AthDate                                map[string]string  `json:"ath_date"`
	Atl                                    map[string]float64 `json:"atl"`
	AtlChangePercentage                    map[string]float64 `json:"atl_change_percentage"`
	AtlDate                                map[string]string  `json:"atl_date"`
	MarketCap                              map[string]float64 `json:"market_cap"`
	MarketCapRank                          int32              `json:"market_cap_rank"`
	FullyDilutedValuation                  map[string]float64 `json:"fully_diluted_valuation"`
	TotalVolume                            map[string]float64 `json:"total_volume"`
	High24h                                map[string]float64 `json:"high_24h"`
	Low24h                                 map[string]float64 `json:"low_24h"`
	PriceChange24h                         float64            `json:"price_change_24h"`
	PriceChangePercentage24h               float64            `json:"price_change_percentage_24h"`
	PriceChangePercentage7d                float64            `json:"price_change_percentage_7d"`
	PriceChangePercentage14d               float64            `json:"price_change_percentage_14d"`
	PriceChangePercentage30d               float64            `json:"price_change_percentage_30d"`
	PriceChangePercentage60d               float64            `json:"price_change_percentage_60d"`
	PriceChangePercentage200d              float64            `json:"price_change_percentage_200d"`
	PriceChangePercentage1y                float64            `json:"price_change_percentage_1y"`
	MarketCapChange24h                     float64            `json:"market_cap_change_24h"`
	MarketCapChangePercentage24h           float64            `json:"market_cap_change_percentage_24h"`
	PriceChange24hInCurrency               map[string]float64 `json:"price_change_24h_in_currency"`
	PriceChangePercentage1hInCurrency      map[string]float64 `json:"price_change_percentage_1h_in_currency"`
	PriceChangePercentage24hInCurrency     map[string]float64 `json:"price_change_percentage_24h_in_currency"`
	PriceChangePercentage7dInCurrency      map[string]float64 `json:"price_change_percentage_7d_in_currency"`
	PriceChangePercentage14dInCurrency     map[string]float64 `json:"price_change_percentage_14d_in_currency"`
	PriceChangePercentage30dInCurrency     map[string]float64 `json:"price_change_percentage_30d_in_currency"`
	PriceChangePercentage60dInCurrency     map[string]float64 `json:"price_change_percentage_60d_in_currency"`
	PriceChangePercentage200dInCurrency    map[string]float64 `json:"price_change_percentage_200d_in_currency"`
	PriceChangePercentage1yInCurrency      map[string]float64 `json:"price_change_percentage_1y_in_currency"`
	MarketCapChange24hInCurrency           map[string]float64 `json:"market_cap_change_24h_in_currency"`
	MarketCapChangePercentage24hInCurrency map[string]float64 `json:"market_cap_change_percentage_24h_in_currency"`
	TotalSupply                            float64            `json:"total_supply"`
	MaxSupply                              float64            `json:"max_supply"`
	CirculatingSupply                      float64            `json:"circulating_supply"`
	LastUpdated                            string             `json:"last_updated"`
}

type CommunityData struct {
	FacebookLikes            interface{} `json:"facebook_likes"`
	TwitterFollowers         int         `json:"twitter_followers"`
	RedditAveragePosts48h    float64     `json:"reddit_average_posts_48h"`
	RedditAverageComments48h float64     `json:"reddit_average_comments_48h"`
	RedditSubscribers        int         `json:"reddit_subscribers"`
	RedditAccountsActive48h  int         `json:"reddit_accounts_active_48h"`
	TelegramChannelUserCount interface{} `json:"telegram_channel_user_count"`
}

type DeveloperData struct {
	Forks              int `json:"forks"`
	Stars              int `json:"stars"`
	Subscribers        int `json:"subscribers"`
	TotalIssues        int `json:"total_issues"`
	ClosedIssues       int `json:"closed_issues"`
	PullRequestsMerged int `json:"pull_requests_merged"`
	PullRequestContributors int `json:"pull_request_contributors"`
	CodeAdditionsDeletions4Weeks struct {
		Additions int `json:"additions"`
		Deletions int `json:"deletions"`
	} `json:"code_additions_deletions_4_weeks"`
	CommitCount4Weeks int `json:"commit_count_4_weeks"`
}

// GetCoinData gets current data (name, price, market, ... including exchange tickers) for a coin
func (c *Client) GetCoinData(ctx context.Context, params *CoinDataRequest) (*CoinDataResponse, *request.Record, error) {
	if params.ID == "" {
		return nil, nil, fmt.Errorf("id is required")
	}

	u, _ := url.Parse(fmt.Sprintf(BaseURL+EndpointCoinData, params.ID))
	q := u.Query()
	
	q.Add(ParamLocalization, fmt.Sprintf("%t", params.Localization))
	q.Add(ParamTickers, fmt.Sprintf("%t", params.Tickers))
	q.Add(ParamMarketData, fmt.Sprintf("%t", params.MarketData))
	q.Add(ParamCommunityData, fmt.Sprintf("%t", params.CommunityData))
	q.Add(ParamDeveloperData, fmt.Sprintf("%t", params.DeveloperData))
	q.Add(ParamSparkline, fmt.Sprintf("%t", params.Sparkline))

	req := request.Request{
		Method: "GET",
		URL:    u.String() + "?" + q.Encode(),
	}

	resp, record, err := c.http.Do(ctx, req)
	if err != nil {
		return nil, record, err
	}

	if resp.StatusCode != 200 {
		return nil, record, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result CoinDataResponse
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return nil, record, err
	}

	return &result, record, nil
}
