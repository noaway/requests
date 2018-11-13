package requests

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
	// _ "golang.org/x/net/proxy"
)

// OptHandle func
type OptHandle func(*Option)

// Option struct
type Option struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	Cookie           http.CookieJar

	ctx   context.Context
	proxy string
}

func defaultOption() Option {
	cookie, _ := cookiejar.New(nil)
	return Option{
		ConnectTimeout:   5 * time.Second,
		ReadWriteTimeout: 5 * time.Second,
		Cookie:           cookie,
		ctx:              context.Background(),
	}
}

// SetProxy func
func SetProxy(proxy string) OptHandle {
	return func(opt *Option) {
		opt.proxy = proxy
	}
}

// SetContext func
func SetContext(ctx context.Context) OptHandle {
	return func(opt *Option) {
		opt.ctx = ctx
	}
}

// HTTPClient struct
type HTTPClient struct {
	url string

	req  *http.Request
	resp *http.Response
	opt  Option

	body []byte
}

// Get func
func (c *HTTPClient) Get(url string, opts ...OptHandle) *HTTPClient {
	return c.newHTTPClient(http.MethodGet, url, opts...)
}

// Post func
func (c *HTTPClient) Post(url string, opts ...OptHandle) *HTTPClient {
	return c.newHTTPClient(http.MethodPost, url, opts...)
}

// Put func
func (c *HTTPClient) Put(url string, opts ...OptHandle) *HTTPClient {
	return c.newHTTPClient(http.MethodPut, url, opts...)
}

// Delete func
func (c *HTTPClient) Delete(url string, opts ...OptHandle) *HTTPClient {
	return c.newHTTPClient(http.MethodDelete, url, opts...)
}

// Head func
func (c *HTTPClient) Head(url string, opts ...OptHandle) *HTTPClient {
	return c.newHTTPClient(http.MethodHead, url, opts...)
}

// SetOpt func
func (c *HTTPClient) SetOpt(opt *Option) *HTTPClient {
	c.opt = *opt
	return c
}

// newHTTPClient func
func (c *HTTPClient) newHTTPClient(method, url string, opts ...OptHandle) *HTTPClient {
	o := defaultOption()
	for i := range opts {
		opts[i](&o)
	}

	return &HTTPClient{
		opt: o,
		url: url,
		req: &http.Request{
			Method: method,
		},
	}
}

func (c *HTTPClient) do() (*http.Response, error) {
	up, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}
	c.req.URL = up

	var transport http.RoundTripper
	if c.opt.proxy != "" {
		transport = &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				return url.Parse(c.opt.proxy)
			},
		}
	}
	httpClient := http.Client{
		Transport: transport,
	}
	return httpClient.Do(c.req)
}

// Bytes func
func (c *HTTPClient) Bytes() ([]byte, error) {
	resp, err := c.do()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// String func
func (c *HTTPClient) String() (string, error) {
	data, err := c.Bytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}
