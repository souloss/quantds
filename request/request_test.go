package request

import (
	"testing"
)

func TestNewRecord(t *testing.T) {
	r := NewRecord()
	if r == nil {
		t.Fatal("NewRecord() returned nil")
	}
	if r.ID == "" {
		t.Error("Record.ID is empty")
	}
	if r.Tags == nil {
		t.Error("Record.Tags is nil")
	}
	if r.StartTime.IsZero() {
		t.Error("Record.StartTime is zero")
	}
}

func TestRecord_IsSuccess(t *testing.T) {
	tests := []struct {
		name        string
		record      *Record
		wantSuccess bool
	}{
		{
			name: "200 status no error",
			record: &Record{
				Response: Response{StatusCode: 200},
			},
			wantSuccess: true,
		},
		{
			name: "201 status no error",
			record: &Record{
				Response: Response{StatusCode: 201},
			},
			wantSuccess: true,
		},
		{
			name: "299 status no error",
			record: &Record{
				Response: Response{StatusCode: 299},
			},
			wantSuccess: true,
		},
		{
			name: "200 status with error",
			record: &Record{
				Response: Response{StatusCode: 200},
				Error:    &RequestError{Type: ErrorTypeNetwork},
			},
			wantSuccess: false,
		},
		{
			name: "400 status",
			record: &Record{
				Response: Response{StatusCode: 400},
			},
			wantSuccess: false,
		},
		{
			name: "500 status",
			record: &Record{
				Response: Response{StatusCode: 500},
			},
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.record.IsSuccess(); got != tt.wantSuccess {
				t.Errorf("Record.IsSuccess() = %v, want %v", got, tt.wantSuccess)
			}
		})
	}
}

func TestRecord_IsError(t *testing.T) {
	tests := []struct {
		name      string
		record    *Record
		wantError bool
	}{
		{
			name: "200 status no error",
			record: &Record{
				Response: Response{StatusCode: 200},
			},
			wantError: false,
		},
		{
			name: "400 status",
			record: &Record{
				Response: Response{StatusCode: 400},
			},
			wantError: true,
		},
		{
			name: "500 status",
			record: &Record{
				Response: Response{StatusCode: 500},
			},
			wantError: true,
		},
		{
			name: "200 status with error",
			record: &Record{
				Response: Response{StatusCode: 200},
				Error:    &RequestError{Type: ErrorTypeNetwork},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.record.IsError(); got != tt.wantError {
				t.Errorf("Record.IsError() = %v, want %v", got, tt.wantError)
			}
		})
	}
}

func TestRecordIDUniqueness(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		r := NewRecord()
		if ids[r.ID] {
			t.Errorf("Duplicate ID generated: %s", r.ID)
		}
		ids[r.ID] = true
	}
}

func TestRecord_ToJSON(t *testing.T) {
	r := &Record{
		ID:    "test-id",
		Error: nil,
		Request: Request{
			Method: "GET",
			URL:    "https://example.com",
		},
		Response: Response{
			StatusCode: 200,
			Body:       []byte("test"),
		},
	}

	data, err := r.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("ToJSON() returned empty data")
	}
}
