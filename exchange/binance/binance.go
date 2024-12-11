package binance

import (
	"context"

	"github.com/go-gotop/gotop/exchange"
	"github.com/go-gotop/gotop/types"
)

// NewBinanceExchange 创建 BinanceExchange
func NewBinanceExchange(opts ...Option) *BinanceExchange {
	o := &options{}
	applyOptions(o, opts...)
	return &BinanceExchange{opts: o}
}

// BinanceExchange 是 Binance 交易所的实现
type BinanceExchange struct {
	opts *options
}

// Name 交易所名称
func (e *BinanceExchange) Name() string {
	return types.BinanceExchange
}

// Assets 交易所支持的资产
func (e *BinanceExchange) Assets(ctx context.Context) ([]types.Asset, error) {
	return nil, nil
}

// CreateOrder 创建订单
func (e *BinanceExchange) CreateOrder(ctx context.Context, req *exchange.CreateOrderRequest) (*exchange.CreateOrderResponse, error) {
	return nil, nil
}

