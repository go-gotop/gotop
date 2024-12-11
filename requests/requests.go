package requests

import (
    "context"
    "encoding/json"
    "errors"
    "io"
    "net/http"

    "strings"
    "time"
)

// HttpClient 定义通用的HTTP调用接口
type HttpClient interface {
	// DoRequest 执行HTTP请求
    DoRequest(ctx context.Context, method, url string, params interface{}) ([]byte, error)
}

// 策略类型定义
// RequestSigner 请求签名
type RequestSigner func(req *http.Request, params interface{}) error
// ErrorHandler 错误处理
type ErrorHandler func(statusCode int, body []byte) error
// HeaderSetter 设置请求头
type HeaderSetter func(req *http.Request, params interface{}) error
// RequestEncoder 请求编码
type RequestEncoder func(method string, urlStr string, params interface{}) (*http.Request, error)

// DefaultJSONRequestEncoder 是默认的JSON编码器
func DefaultJSONRequestEncoder(method, urlStr string, params interface{}) (*http.Request, error) {
    var body io.Reader
    if params != nil {
        data, err := json.Marshal(params)
        if err != nil {
            return nil, err
        }
        body = strings.NewReader(string(data))
    }

    req, err := http.NewRequest(method, urlStr, body)
    if err != nil {
        return nil, err
    }

    if method == http.MethodPost || method == http.MethodPut {
        req.Header.Set("Content-Type", "application/json")
    }

    return req, nil
}

type httpClientImpl struct {
    client         *http.Client
    signer         RequestSigner
    errorHandler   ErrorHandler
    headerSetter   HeaderSetter
    requestEncoder RequestEncoder
}

// NewHttpClient 创建一个可定制化的HttpClient
func NewHttpClient(opts ...HttpClientOption) HttpClient {
    c := &httpClientImpl{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
        requestEncoder: DefaultJSONRequestEncoder,
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}

// DoRequest 执行HTTP请求
func (c *httpClientImpl) DoRequest(ctx context.Context, method, urlStr string, params interface{}) ([]byte, error) {
    req, err := c.requestEncoder(method, urlStr, params)
    if err != nil {
        return nil, err
    }
    req = req.WithContext(ctx)

    if c.headerSetter != nil {
        if err := c.headerSetter(req, params); err != nil {
            return nil, err
        }
    }

    if c.signer != nil {
        if err := c.signer(req, params); err != nil {
            return nil, err
        }
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode < 200 || resp.StatusCode > 299 {
        if c.errorHandler != nil {
            return nil, c.errorHandler(resp.StatusCode, body)
        }
        return nil, errors.New("http error: " + resp.Status)
    }

    return body, nil
}