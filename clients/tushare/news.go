package tushare

import (
	"context"

	"github.com/souloss/quantds/request"
)

// NewsParams 新闻快讯查询参数。
// Tushare API: news
// 获取主流新闻网站的财经资讯。
type NewsParams struct {
	StartDate string // 开始日期 (YYYY-MM-DD HH:MM:SS)
	EndDate   string // 结束日期 (YYYY-MM-DD HH:MM:SS)
	Src       string // 新闻来源 (sina, wallstreetcn, 10jqka, eastmoney, yicai)
}

func (p *NewsParams) ToMap() map[string]string {
	m := make(map[string]string)
	if p == nil {
		return m
	}
	if p.StartDate != "" {
		m["start_date"] = p.StartDate
	}
	if p.EndDate != "" {
		m["end_date"] = p.EndDate
	}
	if p.Src != "" {
		m["src"] = p.Src
	}
	return m
}

// NewsRow 新闻快讯数据行。
type NewsRow struct {
	Datetime string // 发布时间
	Content  string // 新闻内容
	Title    string // 新闻标题
	Channels string // 分类
}

// GetNews 获取新闻快讯。
func (c *Client) GetNews(ctx context.Context, params *NewsParams) ([]NewsRow, *request.Record, error) {
	data, record, err := c.post(ctx, "news", params.ToMap(), "datetime,content,title,channels")
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]NewsRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, NewsRow{
			Datetime: getStr(idx, item, "datetime"),
			Content:  getStr(idx, item, "content"),
			Title:    getStr(idx, item, "title"),
			Channels: getStr(idx, item, "channels"),
		})
	}

	return rows, record, nil
}

// CctvNewsParams 新闻联播参数
type CctvNewsParams struct {
	Date string // 日期 YYYYMMDD
}

func (p *CctvNewsParams) ToMap() map[string]string {
	m := make(map[string]string)
	if p == nil {
		return m
	}
	if p.Date != "" {
		m["date"] = p.Date
	}
	return m
}

// CctvNewsRow 新闻联播数据行
type CctvNewsRow struct {
	Date    string // 日期
	Title   string // 标题
	Content string // 内容
}

// GetCctvNews 获取新闻联播
func (c *Client) GetCctvNews(ctx context.Context, params *CctvNewsParams) ([]CctvNewsRow, *request.Record, error) {
	data, record, err := c.post(ctx, "cctv_news", params.ToMap(), "date,title,content")
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]CctvNewsRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, CctvNewsRow{
			Date:    getStr(idx, item, "date"),
			Title:   getStr(idx, item, "title"),
			Content: getStr(idx, item, "content"),
		})
	}

	return rows, record, nil
}

// AnnouncementParams 上市公司公告查询参数
// Tushare API: disclosure
type AnnouncementParams struct {
	TsCode    string // 股票代码
	AnnDate   string // 公告日期 (YYYYMMDD)
	StartDate string // 开始日期 (YYYYMMDD)
	EndDate   string // 结束日期 (YYYYMMDD)
}

func (p *AnnouncementParams) ToMap() map[string]string {
	m := make(map[string]string)
	if p == nil {
		return m
	}
	if p.TsCode != "" {
		m["ts_code"] = p.TsCode
	}
	if p.AnnDate != "" {
		m["ann_date"] = p.AnnDate
	}
	if p.StartDate != "" {
		m["start_date"] = p.StartDate
	}
	if p.EndDate != "" {
		m["end_date"] = p.EndDate
	}
	return m
}

// AnnouncementRow 上市公司公告数据行
type AnnouncementRow struct {
	TsCode  string // 股票代码
	AnnDate string // 公告日期
	Title   string // 标题
	Content string // 内容 (Tushare disclosure usually doesn't return full content, but we can check fields)
}

// GetAnnouncements 获取上市公司公告
func (c *Client) GetAnnouncements(ctx context.Context, params *AnnouncementParams) ([]AnnouncementRow, *request.Record, error) {
	// 字段: ts_code,ann_date,title,content
	data, record, err := c.post(ctx, "disclosure", params.ToMap(), "ts_code,ann_date,title,content")
	if err != nil {
		return nil, record, err
	}

	idx := fieldIndex(data.Fields)
	rows := make([]AnnouncementRow, 0, len(data.Items))
	for _, item := range data.Items {
		rows = append(rows, AnnouncementRow{
			TsCode:  getStr(idx, item, "ts_code"),
			AnnDate: getStr(idx, item, "ann_date"),
			Title:   getStr(idx, item, "title"),
			Content: getStr(idx, item, "content"),
		})
	}

	return rows, record, nil
}
