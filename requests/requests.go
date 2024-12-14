package requests

import (
	"bytes"
	"errors"
	"net/http"
)

// AuthInfo 包含每次请求所需的 API Key 和 Secret Key 和 Passphrase。
type AuthInfo struct {
	APIKey     string
	SecretKey  string
	Passphrase string
}

// Request 表示一个完整的请求，包括方法、URL、业务参数以及鉴权信息。
type Request struct {
	Method string
	URL    string
	Params map[string]interface{}
	Auth   *AuthInfo
}

// PreparedRequest 代表一个已经构建好的HTTP请求参数，包含方法、完整URL、请求头、以及请求体。
// ExchangeAdapter 将直接返回此对象，以供 HttpClient 直接发送。
type PreparedRequest struct {
	Method  string
	URL     string
	Headers http.Header
	Body    []byte
}

// ExchangeAdapter 接口由各交易所自行实现，用于构建该交易所特定要求的请求。
// 用户在业务代码中为每个交易所提供自己的Adapter实现。
type ExchangeAdapter interface {
	// BuildRequest 根据输入参数构建一个完整的 PreparedRequest。
	// 内部完成：URL拼接、查询参数处理、Headers构建、签名注入、请求体序列化等所有定制逻辑。
	BuildRequest(req *Request) (*PreparedRequest, error)
}

// RequestClient 定义了通用的HTTP客户端接口。
// 只负责调用 ExchangeAdapter 构建请求，并通过 net/http.Client 实际发送请求。
type RequestClient interface {
	// DoRequest 使用注入的Adapter构建请求并发送。
	// 返回的 *http.Response 由调用者自行读取和处理。
	DoRequest(req *Request) (*http.Response, error)

	// SetAdapter 注入适配器，不同交易所使用不同的Adapter实现。
	SetAdapter(adapter ExchangeAdapter)

	// 可选地设置HTTP客户端（如超时、代理等）
	SetHTTPClient(client *http.Client)
}

// client 是 RequestClient 的默认实现，使用默认的http.Client。
type client struct {
	adapter    ExchangeAdapter
	httpClient *http.Client
}

// NewClient 创建一个新的HttpClient实例，并使用默认的http.Client。
func NewClient() RequestClient {
	return &client{
		httpClient: http.DefaultClient,
	}
}

// SetAdapter 注入适配器，不同交易所使用不同的Adapter实现。
func (c *client) SetAdapter(adapter ExchangeAdapter) {
	c.adapter = adapter
}

// SetHTTPClient 可选地设置HTTP客户端（如超时、代理等）
func (c *client) SetHTTPClient(hc *http.Client) {
	c.httpClient = hc
}

// DoRequest 使用已设置的Adapter构建请求，然后通过http.Client发送请求。
func (c *client) DoRequest(r *Request) (*http.Response, error) {
	if c.adapter == nil {
		return nil, errors.New("no adapter set")
	}

	prepared, err := c.adapter.BuildRequest(r)
	if err != nil {
		return nil, err
	}

	// 构建http.Request
	req, err := http.NewRequest(prepared.Method, prepared.URL, bytes.NewReader(prepared.Body))
	if err != nil {
		return nil, err
	}

	if prepared.Headers != nil {
		req.Header = prepared.Headers
	}

	// 通过http.Client发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
