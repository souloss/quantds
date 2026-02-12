package domain

import "testing"

func TestParseSymbol(t *testing.T) {
	tests := []struct {
		symbol       string
		wantCode     string
		wantExchange Exchange
		wantOK       bool
	}{
		{"000001.SZ", "000001", ExchangeSZ, true},
		{"600001.SH", "600001", ExchangeSH, true},
		{"430001.BJ", "430001", ExchangeBJ, true},
		{"SH600001", "600001", ExchangeSH, true},
		{"SZ000001", "000001", ExchangeSZ, true},
		{"BJ430001", "430001", ExchangeBJ, true},
		{"", "", "", false},
		{"600001.XX", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.symbol, func(t *testing.T) {
			code, exchange, ok := ParseSymbol(tt.symbol)
			if ok != tt.wantOK {
				t.Errorf("ParseSymbol() ok = %v, want %v", ok, tt.wantOK)
			}
			if tt.wantOK {
				if code != tt.wantCode {
					t.Errorf("ParseSymbol() code = %v, want %v", code, tt.wantCode)
				}
				if exchange != tt.wantExchange {
					t.Errorf("ParseSymbol() exchange = %v, want %v", exchange, tt.wantExchange)
				}
			}
		})
	}
}

func TestFormatSymbol(t *testing.T) {
	tests := []struct {
		code     string
		exchange Exchange
		want     string
	}{
		{"000001", ExchangeSZ, "000001.SZ"},
		{"600001", ExchangeSH, "600001.SH"},
		{"430001", ExchangeBJ, "430001.BJ"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := FormatSymbol(tt.code, tt.exchange); got != tt.want {
				t.Errorf("FormatSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseAndFormatRoundTrip(t *testing.T) {
	symbols := []string{"000001.SZ", "600001.SH", "430001.BJ"}
	for _, s := range symbols {
		code, exchange, ok := ParseSymbol(s)
		if !ok {
			t.Errorf("ParseSymbol(%s) failed", s)
			continue
		}
		got := FormatSymbol(code, exchange)
		if got != s {
			t.Errorf("Round trip: %s -> (%s, %s) -> %s", s, code, exchange, got)
		}
	}
}
