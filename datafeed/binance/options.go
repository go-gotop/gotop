package binance

import (
	"log/slog"
)

type options struct {
	// logger 日志记录器
	logger *slog.Logger
}

func applyOptions(opts ...Option) *options {
	o := &options{
		logger: slog.Default(),
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Option 是DataFeed的配置选项
type Option func(o *options)

// WithLogger 设置日志记录器
func WithLogger(logger *slog.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}
