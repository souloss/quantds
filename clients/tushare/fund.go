package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

const (
	APIFundBasic = "fund_basic"
	APIFundNAV   = "fund_nav"
)

const (
	FieldsFundBasic = "ts_code,name,management,custodian,fund_type,found_date,due_date,list_date,issue_date,delist_date,issue_amount,m_fee,c_fee,duration,p_value,min_amount,exp_return,benchmark,status,invest_type,type,trustee,purc_startdate,redm_startdate,market"
	FieldsFundNAV   = "ts_code,ann_date,end_date,unit_nav,accum_nav,accum_div,net_asset,total_netasset,adj_nav"
)

type FundBasicParams struct {
	Market   string // E (场内), O (场外)
	Status   string // D (摘牌), I (发行), L (上市中)
	FundType string
}

type FundBasicRow struct {
	TSCode      string
	Name        string
	Management  string
	Custodian   string
	FundType    string
	FoundDate   string
	ListDate    string
	IssueAmount string
	Status      string
	InvestType  string
	Type        string
	Market      string
}

func (c *Client) GetFundBasic(ctx context.Context, params *FundBasicParams) ([]FundBasicRow, *request.Record, error) {
	m := make(map[string]string)
	if params != nil {
		if params.Market != "" {
			m["market"] = params.Market
		}
		if params.Status != "" {
			m["status"] = params.Status
		}
	}

	data, record, err := c.post(ctx, APIFundBasic, m, FieldsFundBasic)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]FundBasicRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, FundBasicRow{
			TSCode:      getStr(idx, item, "ts_code"),
			Name:        getStr(idx, item, "name"),
			Management:  getStr(idx, item, "management"),
			Custodian:   getStr(idx, item, "custodian"),
			FundType:    getStr(idx, item, "fund_type"),
			FoundDate:   getStr(idx, item, "found_date"),
			ListDate:    getStr(idx, item, "list_date"),
			IssueAmount: getStr(idx, item, "issue_amount"),
			Status:      getStr(idx, item, "status"),
			InvestType:  getStr(idx, item, "invest_type"),
			Type:        getStr(idx, item, "type"),
			Market:      getStr(idx, item, "market"),
		})
	}

	return rows, record, nil
}

type FundNAVParams struct {
	TSCode string
	Market string
}

type FundNAVRow struct {
	TSCode        string
	AnnDate       string
	EndDate       string
	UnitNAV       float64
	AccumNAV      float64
	AccumDiv      float64
	NetAsset      float64
	TotalNetAsset float64
	AdjNAV        float64
}

func (c *Client) GetFundNAV(ctx context.Context, params *FundNAVParams) ([]FundNAVRow, *request.Record, error) {
	m := make(map[string]string)
	if params.TSCode != "" {
		m["ts_code"] = params.TSCode
	}
	if params.Market != "" {
		m["market"] = params.Market
	}

	data, record, err := c.post(ctx, APIFundNAV, m, FieldsFundNAV)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]FundNAVRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, FundNAVRow{
			TSCode:        getStr(idx, item, "ts_code"),
			AnnDate:       getStr(idx, item, "ann_date"),
			EndDate:       getStr(idx, item, "end_date"),
			UnitNAV:       getFlt(idx, item, "unit_nav"),
			AccumNAV:      getFlt(idx, item, "accum_nav"),
			AccumDiv:      getFlt(idx, item, "accum_div"),
			NetAsset:      getFlt(idx, item, "net_asset"),
			TotalNetAsset: getFlt(idx, item, "total_netasset"),
			AdjNAV:        getFlt(idx, item, "adj_nav"),
		})
	}

	return rows, record, nil
}
