package binance

import "github.com/go-gotop/gotop/requests"

// options 是 BinanceExchange 的配置选项
type options struct {
	// apiKey 是 Binance API 的 API Key
	apiKey    string
	// secretKey 是 Binance API 的 Secret Key
	secretKey string
	// client 是 Binance API 的 HTTP 客户端
	client requests.HttpClient
}

// Option 是 BinanceExchange 的配置选项
type Option func(o *options)

// WithApiKey 设置 API Key
func WithApiKey(apiKey string) Option {
	return func(o *options) {
		o.apiKey = apiKey
	}
}

// WithSecretKey 设置 Secret Key
func WithSecretKey(secretKey string) Option {
	return func(o *options) {
		o.secretKey = secretKey
	}
}

// WithClient 设置 HTTP 客户端
func WithClient(client requests.HttpClient) Option {
	return func(o *options) {
		o.client = client
	}
}

// applyOptions 应用配置选项
func applyOptions(o *options, opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}
