package types

import "github.com/shopspring/decimal"

const (
	// BinanceExchange 币安
	BinanceExchange = "BINANCE"
	// HuobiExchange 火币
	HuobiExchange = "HUOBI"
	// OkxExchange OKX
	OkxExchange = "OKX"
	// CoinBaseExchange CoinBase
	CoinBaseExchange = "COINBASE"
	// MockExchange 模拟
	MockExchange = "MOCK"

	// Maker 挂单
	Maker = "MAKER"
	// Taker 吃单
	Taker = "TAKER"

	// ByUser 用户
	ByUser = "USER"
	// BySystem 系统
	BySystem = "SYSTEM"
)

// MarketType 市场类型: 1-MarketTypeSpot, 2-MarketTypeFutures, 3-MarketTypeOptions
type MarketType int

// String 返回字符串表示
func (m MarketType) String() string {
	switch m {
	case MarketTypeSpot:
		return "SPOT"
	case MarketTypeFutures:
		return "FUTURES"
	case MarketTypeOptions:
		return "OPTIONS"
	}
	return ""
}

const (
	// MarketTypeSpot 现货市场
	MarketTypeSpot MarketType = iota + 1
	// MarketTypeFutures 期货市场
	MarketTypeFutures
	// MarketTypeOptions 期权市场
	MarketTypeOptions
)

// TradeEvent 成交事件
type TradeEvent struct {
	// ID 事件ID, 自增id由0开始
	// 如若需要用到该ID, 则需要自行实现
	ID uint64
	// Timestamp 事件时间
	Timestamp int64
	// Symbol 交易对
	Symbol string
	// Exchange 交易所
	Exchange string
	// Size 成交数量
	Size decimal.Decimal
	// Price 成交价格
	Price decimal.Decimal
	// Side 成交方向
	Side SideType
	// Type 市场类型
	Type MarketType
}

type Symbol struct {

}

type Asset struct {
	
}