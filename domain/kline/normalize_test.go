package kline

import (
	"math"
	"testing"
)

func TestBar_Normalize(t *testing.T) {
	tests := []struct {
		name string
		bar  Bar
		want Bar
	}{
		{
			name: "fills Change from Open/Close",
			bar:  Bar{Open: 10.0, Close: 11.0, Volume: 1000},
			want: Bar{Open: 10.0, Close: 11.0, Volume: 1000, Change: 1.0, ChangeRate: 10.0, Turnover: 11000.0},
		},
		{
			name: "respects existing Change",
			bar:  Bar{Open: 10.0, Close: 11.0, Change: 0.5, Volume: 1000},
			want: Bar{Open: 10.0, Close: 11.0, Change: 0.5, ChangeRate: 10.0, Turnover: 11000.0},
		},
		{
			name: "respects existing ChangeRate",
			bar:  Bar{Open: 10.0, Close: 11.0, ChangeRate: 5.0, Volume: 1000},
			want: Bar{Open: 10.0, Close: 11.0, Change: 1.0, ChangeRate: 5.0, Turnover: 11000.0},
		},
		{
			name: "respects existing Turnover",
			bar:  Bar{Open: 10.0, Close: 11.0, Turnover: 5000.0, Volume: 1000},
			want: Bar{Open: 10.0, Close: 11.0, Turnover: 5000.0, Change: 1.0, ChangeRate: 10.0},
		},
		{
			name: "skips when Open is zero",
			bar:  Bar{Close: 11.0, Volume: 1000},
			want: Bar{Close: 11.0, Volume: 1000},
		},
		{
			name: "skips Turnover when Volume is zero",
			bar:  Bar{Open: 10.0, Close: 11.0},
			want: Bar{Open: 10.0, Close: 11.0, Change: 1.0, ChangeRate: 10.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.bar.Normalize()
			if tt.bar.Change != tt.want.Change {
				t.Errorf("Change = %v, want %v", tt.bar.Change, tt.want.Change)
			}
			if math.Abs(tt.bar.ChangeRate-tt.want.ChangeRate) > 0.001 {
				t.Errorf("ChangeRate = %v, want %v", tt.bar.ChangeRate, tt.want.ChangeRate)
			}
			if math.Abs(tt.bar.Turnover-tt.want.Turnover) > 0.001 {
				t.Errorf("Turnover = %v, want %v", tt.bar.Turnover, tt.want.Turnover)
			}
		})
	}
}

func TestResponse_Normalize(t *testing.T) {
	resp := Response{
		Bars: []Bar{
			{Open: 10.0, Close: 11.0, Volume: 1000},
			{Open: 20.0, Close: 22.0, Volume: 500},
		},
	}
	resp.Normalize()

	if resp.Bars[0].Change != 1.0 {
		t.Errorf("Bars[0].Change = %v, want 1.0", resp.Bars[0].Change)
	}
	if resp.Bars[1].Change != 2.0 {
		t.Errorf("Bars[1].Change = %v, want 2.0", resp.Bars[1].Change)
	}
}
