package types

import "github.com/shopspring/decimal"

// StrategyStatus 策略状态:
// 1-StrategyStatusRunning, 2-StrategyStatusSuspended
// 3-StrategyStatusStopped, 4-StrategyStatusFinished
// 5-StrategyStatusError
type StrategyStatus int

// String 返回字符串表示
func (s StrategyStatus) String() string {
	switch s {
	case StrategyStatusRunning:
		return "RUNNING"
	case StrategyStatusSuspended:
		return "SUSPENDED"
	case StrategyStatusStopped:
		return "STOPPED"
	case StrategyStatusFinished:
		return "FINISHED"
	case StrategyStatusError:
		return "ERROR"
	}
	return "UNKNOWN"
}

const (
	// StrategyStatusUnknown 未知
	StrategyStatusUnknown StrategyStatus = iota
	// StrategyStatusRunning 运行中
	StrategyStatusRunning
	// StrategyStatusSuspended 挂起
	StrategyStatusSuspended
	// StrategyStatusStopped 停止
	StrategyStatusStopped
	// StrategyStatusFinished 完成
	StrategyStatusFinished
	// StrategyStatusError 错误
	StrategyStatusError
)

// PriceDirection 价格方向: 1-PriceDirectionUp, 2-PriceDirectionDown
type PriceDirection int

// String 返回字符串表示
func (p PriceDirection) String() string {
	switch p {
	case PriceDirectionUp:
		return "UP"
	case PriceDirectionDown:
		return "DOWN"
	}
	return "UNKNOWN"
}

const (
	// PriceDirectionUnknown 未知
	PriceDirectionUnknown PriceDirection = iota
	// PriceDirectionUp 向上
	PriceDirectionUp
	// PriceDirectionDown 向下
	PriceDirectionDown
)

// PositioningLevel 支撑阻力级别: 1-PositioningLevelSupport, 2-PositioningLevelResistance
type PositioningLevel int

// String 返回字符串表示
func (p PositioningLevel) String() string {
	switch p {
	case PositioningLevelSupport:
		return "SUPPORT"
	case PositioningLevelResistance:
		return "RESISTANCE"
	}
	return "UNKNOWN"
}

const (
	// PositioningLevelUnknown 未知
	PositioningLevelUnknown PositioningLevel = iota
	// PositioningLevelSupport 支撑
	PositioningLevelSupport
	// PositioningLevelResistance 阻力
	PositioningLevelResistance
)

// StrategySignal 策略信号
type StrategySignal struct {
	// Price 价格
	Price decimal.Decimal
	// Size 头寸大小
	Size decimal.Decimal
	// Side 方向: 1-SideBuy, 2-SideSell
	Side SideType
	// OrderType 订单类型: 1-OrderTypeMarket, 2-OrderTypeLimit
	OrderType OrderType
	// PositionSide 持仓方向: 1-PositionSideLong, 2-PositionSideShort
	PositionSide PositionSide
}

// PricePoint 价格点
type PricePoint struct {
	// Price 价格
	Price decimal.Decimal
	// ID 全局索引起始为0，递增
	ID uint64
	// Timestamp 交易时间
	Timestamp int64
	// Direction 价格方向: 1-PriceDirectionUp, 2-PriceDirectionDown
	Direction PriceDirection
}

// RangeExtremum 区间极值
type RangeExtremum struct {
	// PeakPrice 最高价格点
	PeakPrice PricePoint
	// ValleyPrice 最低价格点
	ValleyPrice PricePoint
}

// IsTrending 是否处于上升趋势
func (r *RangeExtremum) IsTrending() bool {
	return r.PeakPrice.Timestamp > r.ValleyPrice.Timestamp
}

// IsSideways 是否横盘
func (r *RangeExtremum) IsSideways() bool {
	return r.PeakPrice.Price.Equal(r.ValleyPrice.Price)
}

// PriceRange 价格区间
func(r *RangeExtremum) PriceRange() (PricePoint, PricePoint) {
	if r.IsTrending() {
		return r.ValleyPrice, r.PeakPrice
	}
	return r.PeakPrice, r.ValleyPrice
}

// TradeAggregate 交易聚合数据
type TradeAggregate struct {
	// SellCount 卖单数量
	SellCount uint64
	// BuyCount 买单数量
	BuyCount uint64
	// Timestamp 聚合时间戳
	Timestamp int64
	// PeakPrice 最高价格点
	PeakPrice PricePoint
	// ValleyPrice 最低价格点
	ValleyPrice PricePoint
	// CurrentPrice 当前价格点
	CurrentPrice PricePoint
	// BuyVolume 买单总量
	BuyVolume decimal.Decimal
	// SellVolume 卖单总量
	SellVolume decimal.Decimal
	// BuyAmount 买单总额
	BuyAmount decimal.Decimal
	// SellAmount 卖单总额
	SellAmount decimal.Decimal
}

// IsTrending 是否处于上升趋势
func (t *TradeAggregate) IsTrending() bool {
	if t.PeakPrice.Price.GreaterThan(t.ValleyPrice.Price) && t.PeakPrice.Timestamp > t.ValleyPrice.Timestamp {
		return true
	}
	return false
}

// IsSideways 是否横盘
func (t *TradeAggregate) IsSideways() bool {
	return t.PeakPrice.Price.Equal(t.ValleyPrice.Price)
}

// PriceRange 价格区间
func(t *TradeAggregate) PriceRange() (PricePoint, PricePoint) {
	if t.IsTrending() {
		return t.ValleyPrice, t.PeakPrice
	}
	return t.PeakPrice, t.ValleyPrice
}