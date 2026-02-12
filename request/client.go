package request

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"resty.dev/v3"
)

type Client interface {
	Do(ctx context.Context, req Request) (Response, *Record, error)
	Close()
}

type ClientImpl struct {
	client   *resty.Client
	executor failsafe.Executor[Response]
}

func NewClient(cfg *Config) *ClientImpl {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	client := resty.New()
	client.SetRetryCount(0)

	var policies []failsafe.Policy[Response]

	if cfg.Timeout != nil {
		policies = append(policies, cfg.Timeout)
	}
	if cfg.CircuitBreaker != nil {
		policies = append(policies, cfg.CircuitBreaker)
	}
	if cfg.RetryPolicy != nil {
		policies = append(policies, cfg.RetryPolicy)
	}
	if cfg.RateLimiter != nil {
		policies = append(policies, cfg.RateLimiter)
	}

	executor := failsafe.With[Response](policies...)

	return &ClientImpl{
		client:   client,
		executor: executor,
	}
}

func NewClientWithResty(restyClient *resty.Client, cfg *Config) *ClientImpl {
	if restyClient == nil {
		restyClient = resty.New()
	}
	restyClient.SetRetryCount(0)

	if cfg == nil {
		cfg = DefaultConfig()
	}

	var policies []failsafe.Policy[Response]

	if cfg.Timeout != nil {
		policies = append(policies, cfg.Timeout)
	}
	if cfg.CircuitBreaker != nil {
		policies = append(policies, cfg.CircuitBreaker)
	}
	if cfg.RetryPolicy != nil {
		policies = append(policies, cfg.RetryPolicy)
	}
	if cfg.RateLimiter != nil {
		policies = append(policies, cfg.RateLimiter)
	}

	executor := failsafe.With[Response](policies...)

	return &ClientImpl{
		client:   restyClient,
		executor: executor,
	}
}

func (c *ClientImpl) Do(ctx context.Context, req Request) (Response, *Record, error) {
	record := NewRecord()
	record.Request = req

	var attempt int
	resp, execErr := c.executor.GetWithExecution(func(exec failsafe.Execution[Response]) (Response, error) {
		attempt = exec.Attempts()
		return c.doHTTP(ctx, req)
	})

	record.Attempt = attempt
	record.Duration = time.Since(record.StartTime)
	record.Response = resp

	if execErr != nil {
		callErr := ClassifyError(execErr, resp.StatusCode)
		record.Error = callErr
		return resp, record, callErr
	}

	if resp.StatusCode >= 400 {
		callErr := ClassifyError(nil, resp.StatusCode)
		record.Error = callErr
		return resp, record, callErr
	}

	return resp, record, nil
}

func (c *ClientImpl) doHTTP(ctx context.Context, req Request) (Response, error) {
	r := c.client.R().SetContext(ctx)

	for k, v := range req.Headers {
		r.SetHeader(k, v)
	}

	if len(req.Body) > 0 {
		r.SetBody(req.Body)
	}

	var resp *resty.Response
	var err error

	switch req.Method {
	case "GET":
		resp, err = r.Get(req.URL)
	case "POST":
		resp, err = r.Post(req.URL)
	case "PUT":
		resp, err = r.Put(req.URL)
	case "DELETE":
		resp, err = r.Delete(req.URL)
	case "PATCH":
		resp, err = r.Patch(req.URL)
	default:
		resp, err = r.Get(req.URL)
	}

	if err != nil {
		return Response{}, err
	}

	result := Response{
		StatusCode: resp.StatusCode(),
		Headers:    make(map[string]string),
		Body:       resp.Bytes(),
	}

	for k, v := range resp.Header() {
		if len(v) > 0 {
			result.Headers[k] = v[0]
		}
	}

	return result, nil
}

func (c *ClientImpl) Close() {
	c.client.Close()
}

var _ Client = (*ClientImpl)(nil)

type NoopClient struct{}

func NewNoopClient() *NoopClient {
	return &NoopClient{}
}

func (c *NoopClient) Do(_ context.Context, req Request) (Response, *Record, error) {
	record := NewRecord()
	record.Request = req
	record.Duration = time.Since(record.StartTime)
	return Response{}, record, nil
}

func (c *NoopClient) Close() {}

var _ Client = (*NoopClient)(nil)

func BuildCacheKey(req Request) string {
	h := sha256.New()
	h.Write([]byte(req.Method))
	h.Write([]byte(req.URL))
	h.Write(req.Body)
	return hex.EncodeToString(h.Sum(nil))
}

type CachingClient struct {
	client Client
	cache  map[string]cachedResponse
	ttl    time.Duration
}

type cachedResponse struct {
	data      []byte
	expiresAt time.Time
}

func NewCachingClient(client Client, ttl time.Duration) *CachingClient {
	return &CachingClient{
		client: client,
		cache:  make(map[string]cachedResponse),
		ttl:    ttl,
	}
}

func (c *CachingClient) Do(ctx context.Context, req Request) (Response, *Record, error) {
	if req.Method == "GET" && c.ttl > 0 {
		key := BuildCacheKey(req)
		if cached, ok := c.cache[key]; ok && time.Now().Before(cached.expiresAt) {
			var resp Response
			if err := json.Unmarshal(cached.data, &resp); err == nil {
				record := NewRecord()
				record.Request = req
				record.FromCache = true
				record.Response = resp
				return resp, record, nil
			}
		}

		resp, record, err := c.client.Do(ctx, req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if data, merr := json.Marshal(resp); merr == nil {
				c.cache[key] = cachedResponse{
					data:      data,
					expiresAt: time.Now().Add(c.ttl),
				}
			}
		}
		return resp, record, err
	}

	return c.client.Do(ctx, req)
}

func (c *CachingClient) Close() {
	c.client.Close()
}

func (c *CachingClient) ClearCache() {
	c.cache = make(map[string]cachedResponse)
}
