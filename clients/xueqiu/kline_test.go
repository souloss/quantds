package xueqiu

import (
	"testing"
)

func TestClient_GetKline(t *testing.T) {
	t.Skip("Xueqiu requires authentication, skipping")
}

func TestToXueqiuSymbol(t *testing.T) {
	tests := []struct {
		symbol  string
		want    string
		wantErr bool
	}{
		{"600001.SH", "SH600001", false},
		{"000001.SZ", "SZ000001", false},
		{"INVALID", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			got, err := toXueqiuSymbol(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("toXueqiuSymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("toXueqiuSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}
