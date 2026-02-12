package manager

import (
	"sync/atomic"
	"time"

	"github.com/souloss/quantds/request"
)

type RequestTrace struct {
	FetchID   string
	Provider  string
	Requests  []*request.Record
	TotalTime time.Duration
	StartTime time.Time
}

func NewRequestTrace(provider string) *RequestTrace {
	return &RequestTrace{
		FetchID:   generateFetchID(),
		Provider:  provider,
		Requests:  make([]*request.Record, 0),
		StartTime: time.Now(),
	}
}

func (t *RequestTrace) AddRequest(record *request.Record) {
	if record != nil {
		t.Requests = append(t.Requests, record)
	}
}

func (t *RequestTrace) TotalRequests() int {
	return len(t.Requests)
}

func (t *RequestTrace) SuccessRequests() int {
	count := 0
	for _, r := range t.Requests {
		if r.IsSuccess() {
			count++
		}
	}
	return count
}

func (t *RequestTrace) FailedRequests() int {
	count := 0
	for _, r := range t.Requests {
		if r.IsError() {
			count++
		}
	}
	return count
}

func (t *RequestTrace) TotalDuration() time.Duration {
	var total time.Duration
	for _, r := range t.Requests {
		total += r.Duration
	}
	return total
}

func (t *RequestTrace) Finish() {
	t.TotalTime = time.Since(t.StartTime)
}

type fetchIDGenerator struct {
	counter uint64
}

func (g *fetchIDGenerator) Next() string {
	n := atomic.AddUint64(&g.counter, 1)
	return time.Now().Format("20060102150405") + "-" + uint64ToString(n)
}

func uint64ToString(n uint64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

var globalFetchID = &fetchIDGenerator{}

func generateFetchID() string {
	return globalFetchID.Next()
}
