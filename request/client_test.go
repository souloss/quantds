package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/failsafe-go/failsafe-go/timeout"
)

func TestNewClient(t *testing.T) {
	client := NewClient(nil)
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	defer client.Close()
}

func TestClient_Do_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	client := NewClient(DefaultConfig())
	defer client.Close()

	req := Request{
		Method:  "GET",
		URL:     ts.URL,
		Headers: map[string]string{"X-Test": "value"},
	}

	resp, record, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if string(resp.Body) != `{"status":"ok"}` {
		t.Errorf("Body = %s, want %s", resp.Body, `{"status":"ok"}`)
	}

	if !record.IsSuccess() {
		t.Error("Record should be success")
	}

	if record.Attempt != 1 {
		t.Errorf("Attempt = %d, want 1", record.Attempt)
	}
}

func TestClient_Do_Post(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`created`))
	}))
	defer ts.Close()

	client := NewClient(DefaultConfig())
	defer client.Close()

	req := Request{
		Method: "POST",
		URL:    ts.URL,
		Body:   []byte(`{"data":"test"}`),
	}

	resp, record, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	if !record.IsSuccess() {
		t.Error("Record should be success")
	}
}

func TestClient_Do_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := NewClient(DefaultConfig(
		WithRetryPolicy(DefaultRetryPolicy()),
	))
	defer client.Close()

	req := Request{
		Method: "GET",
		URL:    ts.URL,
	}

	_, record, err := client.Do(context.Background(), req)
	if err == nil {
		t.Error("Expected error for 500 status")
	}

	if record.Error == nil {
		t.Error("Record.Error should not be nil")
	}

	if record.Error.Type != ErrorTypeServer {
		t.Errorf("Error.Type = %s, want %s", record.Error.Type, ErrorTypeServer)
	}
}

func TestClient_Do_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := DefaultConfig()
	cfg.Timeout = timeout.New[Response](100 * time.Millisecond)

	client := NewClient(cfg)
	defer client.Close()

	req := Request{
		Method: "GET",
		URL:    ts.URL,
	}

	_, record, err := client.Do(context.Background(), req)
	if err == nil {
		t.Error("Expected timeout error")
	}

	if record.Error == nil {
		t.Error("Record.Error should not be nil for timeout")
	}

	if record.Error.Type != ErrorTypeTimeout {
		t.Errorf("Error.Type = %s, want %s", record.Error.Type, ErrorTypeTimeout)
	}
}

func TestClient_Do_ConnectionError(t *testing.T) {
	client := NewClient(DefaultConfig())
	defer client.Close()

	req := Request{
		Method: "GET",
		URL:    "http://127.0.0.1:99999/nonexistent",
	}

	_, record, err := client.Do(context.Background(), req)
	if err == nil {
		t.Error("Expected connection error")
	}

	if record.Error == nil {
		t.Error("Record.Error should not be nil")
	}
}

func TestNoopClient(t *testing.T) {
	client := NewNoopClient()
	defer client.Close()

	req := Request{
		Method: "GET",
		URL:    "http://example.com",
	}

	resp, record, err := client.Do(context.Background(), req)
	if err != nil {
		t.Errorf("NoopClient.Do() error = %v", err)
	}

	if record.Request.URL != req.URL {
		t.Errorf("Record.Request.URL = %s, want %s", record.Request.URL, req.URL)
	}

	if resp.StatusCode != 0 {
		t.Errorf("Response.StatusCode = %d, want 0", resp.StatusCode)
	}
}

func TestCachingClient(t *testing.T) {
	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`response`))
	}))
	defer ts.Close()

	innerClient := NewClient(DefaultConfig())
	client := NewCachingClient(innerClient, 5*time.Minute)
	defer client.Close()

	req := Request{
		Method: "GET",
		URL:    ts.URL,
	}

	_, _, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("First Do() error = %v", err)
	}
	if callCount != 1 {
		t.Errorf("First call: callCount = %d, want 1", callCount)
	}

	_, record, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Second Do() error = %v", err)
	}
	if callCount != 1 {
		t.Errorf("Second call (cached): callCount = %d, want 1", callCount)
	}
	if !record.FromCache {
		t.Error("Record.FromCache should be true")
	}

	client.ClearCache()
	_, _, err = client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("Third Do() error = %v", err)
	}
	if callCount != 2 {
		t.Errorf("Third call (after clear): callCount = %d, want 2", callCount)
	}
}

func TestBuildCacheKey(t *testing.T) {
	req1 := Request{Method: "GET", URL: "http://example.com", Body: []byte("test")}
	req2 := Request{Method: "GET", URL: "http://example.com", Body: []byte("test")}
	req3 := Request{Method: "POST", URL: "http://example.com", Body: []byte("test")}

	key1 := BuildCacheKey(req1)
	key2 := BuildCacheKey(req2)
	key3 := BuildCacheKey(req3)

	if key1 != key2 {
		t.Error("Same requests should have same cache key")
	}

	if key1 == key3 {
		t.Error("Different methods should have different cache keys")
	}
}
