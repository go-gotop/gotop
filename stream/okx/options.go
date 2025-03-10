package okx

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Option func(*OkxStream)

func applyOptions(b *OkxStream, opts ...Option) {
	for _, opt := range opts {
		opt(b)
	}
}

// WithPingInterval 设置Ping间隔
func WithPingInterval(d time.Duration) Option {
	return func(b *OkxStream) {
		b.pingInterval = d
	}
}

// WithPongWait 设置Pong等待时间
func WithPongWait(d time.Duration) Option {
	return func(b *OkxStream) {
		b.pongWait = d
	}
}

// WithWriteWait 设置写等待时间
func WithWriteWait(d time.Duration) Option {
	return func(b *OkxStream) {
		b.writeWait = d
	}
}

// WithReconnectInterval 设置重连间隔
func WithReconnectInterval(d time.Duration) Option {
	return func(b *OkxStream) {
		b.reconnectInterval = d
	}
}

// WithDialer 设置自定义的dial函数
func WithDialer(dialer func(urlStr string, requestHeader http.Header) (*websocket.Conn, *http.Response, error)) Option {
	return func(b *OkxStream) {
		b.dialer = dialer
	}
}
