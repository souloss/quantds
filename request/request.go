package request

import (
	"encoding/json"
	"sync/atomic"
	"time"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

type Record struct {
	ID        string
	Request   Request
	Response  Response
	StartTime time.Time
	Duration  time.Duration
	Error     *RequestError
	Attempt   int
	FromCache bool
	Tags      map[string]string
}

func (r *Record) IsSuccess() bool {
	return r.Error == nil && r.Response.StatusCode >= 200 && r.Response.StatusCode < 300
}

func (r *Record) IsError() bool {
	return r.Error != nil || r.Response.StatusCode >= 400
}

func (r *Record) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

type recordIDGenerator struct {
	counter uint64
}

func (g *recordIDGenerator) Next() string {
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

var globalRecordID = &recordIDGenerator{}

func NewRecord() *Record {
	return &Record{
		ID:        globalRecordID.Next(),
		StartTime: time.Now(),
		Tags:      make(map[string]string),
	}
}
