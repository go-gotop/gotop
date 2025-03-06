package broker

import (
	"github.com/go-gotop/gotop/types"
	"github.com/shopspring/decimal"
)

// StrategySignalEvent 策略信号事件
type StrategySignalEvent struct {
	// PositionID 仓位ID
	PositionID string
	// ID 交易ID
	TransactionID string
	// AccountID 账户ID
	AccountID string
	// Timestamp 当前时间戳
	Timestamp int64
	// ClientOrderID 自定义客户端订单号
	ClientOrderID string
	// Exchange 交易所
	Exchange string
	// TimeInForce GTC，IOC，FOK，GTX，GTD
	TimeInForce types.TimeInForce
	// SideType BUY，SELL
	Side types.SideType
	// OrderType LIMIT，MARKET
	OrderType types.OrderType
	// PositionSide LONG，SHORT
	PositionSide types.PositionSide
	// MarketType 种类 SPOT, FUTURES, MARGIN
	MarketType types.MarketType
	// Symbol 交易对
	Symbol string
	// Size 头寸数量
	Size decimal.Decimal
	// Price 交易价格
	Price decimal.Decimal
	// CreatedBy 创建者 USER, SYSTEM
	CreatedBy string
}

// OrderResultEvent 订单结果事件
type OrderResultEvent struct {
	// AccountID 账户ID
	AccountID string
	// ID 交易ID
	TransactionID string
	// Exchange 交易所
	Exchange string
	// PositionID 仓位ID
	PositionID string
	// ClientOrderID 自定义客户端订单号
	ClientOrderID string
	// Symbol 交易对
	Symbol string
	// OrderID 交易所订单号
	OrderID string
	// FeeAsset 手续费资产
	FeeAsset string
	// TransactionTime 交易时间
	TransactionTime int64
	// By 是否是挂单方 MAKER, TAKER
	By string
	// CreatedBy 创建者 USER，SYSTEM
	CreatedBy string
	// Instrument 种类 SPOT, FUTURES
	MarketType types.MarketType
	// Status 订单状态: OpeningPosition, HoldingPosition, ClosingPosition, ClosedPosition
	Status types.PositionStatus
	// ExecutionType 本次订单执行类型:NEW, TRADE, CANCELED, REJECTED, EXPIRED
	ExecutionType types.ExecutionType
	// State 当前订单执行类型:NEW, PARTIALLY_FILLED, FILLED, CANCELED, REJECTED, EXPIRED
	State types.OrderState
	// PositionSide LONG，SHORT
	PositionSide types.PositionSide
	// SideType BUY，SELL
	Side types.SideType
	// OrderType LIMIT，MARKET
	Type types.OrderType
	// Volume 原交易数量
	Volume decimal.Decimal
	// Price 交易价格
	Price decimal.Decimal
	// LatestVolume 最新交易数量
	LatestVolume decimal.Decimal
	// FilledVolume 已成交数量
	FilledVolume decimal.Decimal
	// LatestPrice 最新交易价格
	LatestPrice decimal.Decimal
	// FeeCost 手续费
	FeeCost decimal.Decimal
	// FilledQuoteVolume 已成交金额
	FilledQuoteVolume decimal.Decimal
	// LatestQuoteVolume 最新成交金额
	LatestQuoteVolume decimal.Decimal
	// QuoteVolume 交易金额
	QuoteVolume decimal.Decimal
	// AvgPrice 平均成交价格
	AvgPrice decimal.Decimal
}

// FrameErrorEvent 帧错误事件
type FrameErrorEvent struct {
	// Error 错误信息
	Error string
	// PositionID 仓位ID
	PositionID string
	// TransactionID 交易ID
	TransactionID string
	// AccountID 账户ID
	AccountID string
	// Timestamp 当前时间戳
	Timestamp int64
	// ClientOrderID 自定义客户端订单号
	ClientOrderID string
}

// TradeEvent 交易事件
type TradeEvent struct {
	// Timestamp 当前时间戳
	Timestamp int64
	// TradeID 交易ID
	TradeID string
	// Symbol 交易对
	Symbol string
	// Exchange 交易所
	Exchange string
	// Volume 交易数量
	Volume decimal.Decimal
	// Price 交易价格
	Price decimal.Decimal
	// Side 交易方向
	Side types.SideType
	// MarketType 市场类型
	MarketType types.MarketType
}

// KlineEvent K线事件
type KlineEvent struct {
	// Symbol 交易对
	Symbol string
	// OpenTime 开盘时间
	OpenTime int64
	// Open 开盘价
	Open decimal.Decimal
	// High 最高价
	High decimal.Decimal
	// Low 最低价
	Low decimal.Decimal
	// Close 收盘价
	Close decimal.Decimal
	// Volume 成交量
	Volume decimal.Decimal
	// CloseTime 收盘时间
	CloseTime int64
	// QuoteAssetVolume 成交额
	QuoteAssetVolume decimal.Decimal
	// NumberOfTrades 成交笔数
	NumberOfTrades int64
	// TakerBuyBaseAssetVolume 买方成交量
	TakerBuyBaseAssetVolume decimal.Decimal
	// TakerBuyQuoteAssetVolume 买方成交额
	TakerBuyQuoteAssetVolume decimal.Decimal
	// Confirm 0 代表 K 线未完结，1 代表 K 线已完结。
	Confirm string
}
