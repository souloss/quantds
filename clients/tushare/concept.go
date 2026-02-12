package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

// ConceptParams 概念分类查询参数。
// Tushare API: concept
// 获取概念板块列表。
type ConceptParams struct {
	Src string // 来源，默认为ts
}

func (p *ConceptParams) ToMap() map[string]string {
	m := make(map[string]string)
	if p == nil {
		return m
	}
	if p.Src != "" {
		m["src"] = p.Src
	}
	return m
}

// ConceptRow 概念分类数据行。
type ConceptRow struct {
	Code string // 概念代码
	Name string // 概念名称
	Src  string // 来源
}

// ConceptDetailParams 概念明细查询参数。
// Tushare API: concept_detail
// 获取概念下的股票列表。
type ConceptDetailParams struct {
	ID     string // 概念代码
	TSCode string // 股票代码
}

func (p *ConceptDetailParams) ToMap() map[string]string {
	m := make(map[string]string)
	if p == nil {
		return m
	}
	if p.ID != "" {
		m["id"] = p.ID
	}
	if p.TSCode != "" {
		m["ts_code"] = p.TSCode
	}
	return m
}

// ConceptDetailRow 概念明细数据行。
type ConceptDetailRow struct {
	ID          string // 概念代码
	ConceptName string // 概念名称
	TSCode      string // 股票代码
	Name        string // 股票名称
	InDate      string // 纳入日期
	OutDate     string // 剔除日期
}

// GetConcept 获取概念分类列表。
func (c *Client) GetConcept(ctx context.Context, params *ConceptParams) ([]ConceptRow, *request.Record, error) {
	data, record, err := c.post(ctx, APIConcept, params.ToMap(), FieldsConcept)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]ConceptRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, ConceptRow{
			Code: getStr(idx, item, "code"),
			Name: getStr(idx, item, "name"),
			Src:  getStr(idx, item, "src"),
		})
	}

	return rows, record, nil
}

// GetConceptDetail 获取概念下的股票列表。
func (c *Client) GetConceptDetail(ctx context.Context, params *ConceptDetailParams) ([]ConceptDetailRow, *request.Record, error) {
	data, record, err := c.post(ctx, APIConceptDetail, params.ToMap(), FieldsConceptDetail)
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]ConceptDetailRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, ConceptDetailRow{
			ID:          getStr(idx, item, "id"),
			ConceptName: getStr(idx, item, "concept_name"),
			TSCode:      getStr(idx, item, "ts_code"),
			Name:        getStr(idx, item, "name"),
			InDate:      getStr(idx, item, "in_date"),
			OutDate:     getStr(idx, item, "out_date"),
		})
	}

	return rows, record, nil
}
