package spot

import (
	"math"
	"testing"
)

func TestQuote_Normalize(t *testing.T) {
	tests := []struct {
		name  string
		quote Quote
		want  Quote
	}{
		{
			name:  "fills Change from Latest/PreClose",
			quote: Quote{Latest: 11.0, PreClose: 10.0, High: 11.5, Low: 9.5, Volume: 1000},
			want:  Quote{Latest: 11.0, PreClose: 10.0, High: 11.5, Low: 9.5, Volume: 1000, Change: 1.0, ChangeRate: 10.0, Amplitude: 20.0, Turnover: 11000.0},
		},
		{
			name:  "respects existing Change",
			quote: Quote{Latest: 11.0, PreClose: 10.0, Change: 0.5, High: 11.5, Low: 9.5, Volume: 1000},
			want:  Quote{Latest: 11.0, PreClose: 10.0, Change: 0.5, ChangeRate: 10.0, Amplitude: 20.0, Turnover: 11000.0},
		},
		{
			name:  "respects existing ChangeRate",
			quote: Quote{Latest: 11.0, PreClose: 10.0, ChangeRate: 5.0, High: 11.5, Low: 9.5, Volume: 1000},
			want:  Quote{Latest: 11.0, PreClose: 10.0, Change: 1.0, ChangeRate: 5.0, Amplitude: 20.0, Turnover: 11000.0},
		},
		{
			name:  "respects existing Amplitude",
			quote: Quote{Latest: 11.0, PreClose: 10.0, Amplitude: 15.0, High: 11.5, Low: 9.5, Volume: 1000},
			want:  Quote{Latest: 11.0, PreClose: 10.0, Change: 1.0, ChangeRate: 10.0, Amplitude: 15.0, Turnover: 11000.0},
		},
		{
			name:  "skips when PreClose is zero",
			quote: Quote{Latest: 11.0, High: 11.5, Low: 9.5, Volume: 1000},
			want:  Quote{Latest: 11.0, High: 11.5, Low: 9.5, Volume: 1000, Turnover: 11000.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.quote.Normalize()
			if tt.quote.Change != tt.want.Change {
				t.Errorf("Change = %v, want %v", tt.quote.Change, tt.want.Change)
			}
			if math.Abs(tt.quote.ChangeRate-tt.want.ChangeRate) > 0.001 {
				t.Errorf("ChangeRate = %v, want %v", tt.quote.ChangeRate, tt.want.ChangeRate)
			}
			if math.Abs(tt.quote.Amplitude-tt.want.Amplitude) > 0.001 {
				t.Errorf("Amplitude = %v, want %v", tt.quote.Amplitude, tt.want.Amplitude)
			}
			if math.Abs(tt.quote.Turnover-tt.want.Turnover) > 0.001 {
				t.Errorf("Turnover = %v, want %v", tt.quote.Turnover, tt.want.Turnover)
			}
		})
	}
}

func TestResponse_Normalize(t *testing.T) {
	resp := Response{
		Quotes: []Quote{
			{Latest: 11.0, PreClose: 10.0, High: 11.5, Low: 9.5, Volume: 1000},
			{Latest: 22.0, PreClose: 20.0, High: 22.5, Low: 19.5, Volume: 500},
		},
	}
	resp.Normalize()

	if resp.Quotes[0].Change != 1.0 {
		t.Errorf("Quotes[0].Change = %v, want 1.0", resp.Quotes[0].Change)
	}
	if resp.Quotes[1].Change != 2.0 {
		t.Errorf("Quotes[1].Change = %v, want 2.0", resp.Quotes[1].Change)
	}
}
