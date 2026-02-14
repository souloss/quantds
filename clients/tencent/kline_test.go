//go:build integration
// +build integration

package tencent

import (
	"context"
	"testing"
	"time"
)

func TestClient_GetKline(t *testing.T) {
	client := NewClient()
	defer client.Close()

	tests := []struct {
		name   string
		params *KlineParams
	}{
		{
			name: "daily kline for SZ stock",
			params: &KlineParams{
				Symbol: "000001.SZ",
				Period: "day",
				Count:  30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			result, record, err := client.GetKline(ctx, tt.params)
			if err != nil {
				t.Fatalf("GetKline() error = %v", err)
			}

			if record == nil {
				t.Fatal("record is nil")
			}

			if !record.IsSuccess() {
				t.Errorf("record should be success, got error: %v", record.Error)
			}

			if result == nil {
				t.Fatal("result is nil")
			}

			if len(result.Data) == 0 {
				t.Log("Warning: no data returned")
				return
			}

			t.Logf("Got %d bars from %s to %s",
				len(result.Data),
				result.Data[0].Date,
				result.Data[len(result.Data)-1].Date)
		})
	}
}

func TestToTencentSymbol(t *testing.T) {
	tests := []struct {
		symbol  string
		want    string
		wantErr bool
	}{
		{"600001.SH", "sh600001", false},
		{"000001.SZ", "sz000001", false},
		{"INVALID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			got, err := toTencentSymbol(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("toTencentSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toTencentSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}
