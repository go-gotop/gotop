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

type Asset struct {
	// AssetName 资产名称
	AssetName string
	// Exchange 交易所
	Exchange string
	// MarketType 市场类型:
	// 1-MarketTypeSpot, 2-MarketTypeFutures, 3-MarketTypeOptions
	MarketType MarketType
	// Free 可用余额
	Free decimal.Decimal
	// Locked 锁定余额
	Locked decimal.Decimal
}

type Symbol struct {
	// OriginalSymbol 原标的物名称
	OriginalSymbol string
	// UnifiedSymbol 统一标的物名称
	UnifiedSymbol string
	// OriginalAsset 原资产名称
	OriginalAsset string
	// UnifiedAsset 统一资产名称
	UnifiedAsset string
	// 交易所
	Exchange string
	// Type 市场类型:
	// 1-MarketTypeSpot, 2-MarketTypeFutures, 3-MarketTypeOptions
	Type MarketType
	// Status 状态: ENABLED, DISABLED
	Status string
	// MinSize 最小头寸
	MinSize decimal.Decimal
	// MaxSize 最大头寸
	MaxSize decimal.Decimal
	// MinPrice 最小价格
	MinPrice decimal.Decimal
	// MaxPrice 最大价格
	MaxPrice decimal.Decimal
	// PricePrecision 价格精度
	PricePrecision int32
	// SizePrecision 头寸精度
	SizePrecision int32
	// CtVal 合约面值
	CtVal decimal.Decimal
	// CtMult 合约乘数
	CtMult decimal.Decimal
	// ListTime 上线时间
	ListTime int64
	// ExpTime 下线时间
	ExpTime int64
}