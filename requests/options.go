package requests

import (
	"net/http"
	"net/url"
	"time"
)

// HttpClientOption 定义客户端选项函数类型
type HttpClientOption func(*httpClientImpl)

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) HttpClientOption {
	return func(c *httpClientImpl) {
		if c.client == nil {
			c.client = &http.Client{}
		}
		c.client.Timeout = timeout
	}
}

// WithClient 设置自定义的HTTP客户端
func WithClient(client *http.Client) HttpClientOption {
	return func(c *httpClientImpl) {
		c.client = client
	}
}

// WithRequestSigner 设置请求签名器
func WithRequestSigner(signer RequestSigner) HttpClientOption {
	return func(c *httpClientImpl) {
		c.signer = signer
	}
}

// WithErrorHandler 设置错误处理器
func WithErrorHandler(handler ErrorHandler) HttpClientOption {
	return func(c *httpClientImpl) {
		c.errorHandler = handler
	}
}

// WithHeaderSetter 设置请求头设置器
func WithHeaderSetter(setter HeaderSetter) HttpClientOption {
	return func(c *httpClientImpl) {
		c.headerSetter = setter
	}
}

// WithRequestEncoder 设置请求编码器
func WithRequestEncoder(encoder RequestEncoder) HttpClientOption {
	return func(c *httpClientImpl) {
		c.requestEncoder = encoder
	}
}

// WithHTTPClient: 允许直接传入自定义的 http.Client
func WithHTTPClient(client *http.Client) HttpClientOption {
	return func(c *httpClientImpl) {
		c.client = client
	}
}

// WithProxy: 简化代理设置的Option
func WithProxy(proxyURL string) HttpClientOption {
	return func(c *httpClientImpl) {
		// 建立一个自定义 Transport
		transport := &http.Transport{}
		if proxyURL != "" {
			parsed, err := url.Parse(proxyURL)
			if err == nil {
				transport.Proxy = http.ProxyURL(parsed)
			}
		}
		c.client.Transport = transport
	}
}