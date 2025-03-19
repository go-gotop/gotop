package types

import (
	"fmt"
	"strings"
)

// SizeUnit 数量单位: 1-SizeUnitCoin 币, 2-SizeUnitContract 合约, 3-SizeUnitQuote 计价货币
type SizeUnit int

// String 返回字符串表示
func (s SizeUnit) String() string {
	switch s {
	case SizeUnitUnknown:
		return "UNKNOWN"
	case SizeUnitCoin:
		return "COIN"
	case SizeUnitContract:
		return "CONTRACT"
	case SizeUnitQuote:
		return "QUOTE"
	}
	return "UNKNOWN"
}

// IsValid 判断 SizeUnit 是否为已定义的类型
func (s SizeUnit) IsValid() bool {
	switch s {
	case SizeUnitUnknown,
		SizeUnitCoin,
		SizeUnitContract,
		SizeUnitQuote:
		return true
	default:
		return false
	}
}

// ParseSizeUnit 从字符串解析 SizeUnit (不区分大小写)
func ParseSizeUnit(s string) (SizeUnit, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "UNKNOWN":
		return SizeUnitUnknown, nil
	case "COIN":
		return SizeUnitCoin, nil
	case "CONTRACT":
		return SizeUnitContract, nil
	case "QUOTE":
		return SizeUnitQuote, nil
	default:
		return SizeUnitUnknown, fmt.Errorf("unknown size unit: %s", s)
	}
}

const (
	// SizeUnitUnknown 未知
	SizeUnitUnknown SizeUnit = iota
	// SizeUnitCoin 币
	SizeUnitCoin
	// SizeUnitContract 合约
	SizeUnitContract
	// SizeUnitQuote 计价货币
	SizeUnitQuote
)

// PosMod 持仓模式：1-Isolated 逐仓, 2-Cross 全仓
type PosMode int

// String 返回字符串表示
func (p PosMode) String() string {
	switch p {
	case PosModeIsolated:
		return "ISOLATED"
	case PosModeCross:
		return "CROSS"
	}
	return "UNKNOWN"
}

// IsValid 判断 PosMod 是否为已定义的类型
func (p PosMode) IsValid() bool {
	switch p {
	case PosModeIsolated,
		PosModeCross:
		return true
	default:
		return false
	}
}

// ParsePosMode 从字符串解析 PosMode (不区分大小写)
func ParsePosMode(s string) (PosMode, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "ISOLATED":
		return PosModeIsolated, nil
	case "CROSS":
		return PosModeCross, nil
	default:
		return PosModeUnknown, fmt.Errorf("unknown position mode: %s", s)
	}
}

const (
	// PosModeUnknown 未知
	PosModeUnknown PosMode = iota
	// PosModeIsolated 逐仓
	PosModeIsolated
	// PosModeCross 全仓
	PosModeCross
)

// ExecutionType 订单执行类型: 1-ExecutionTypeNew, 2-ExecutionTypeTrade, 3-ExecutionTypeCanceled, 4-ExecutionTypeRejected, 5-ExecutionTypeExpired
type ExecutionType int

// String 返回字符串表示
func (e ExecutionType) String() string {
	switch e {
	case ExecutionTypeNew:
		return "NEW"
	case ExecutionTypeTrade:
		return "TRADE"
	case ExecutionTypeCanceled:
		return "CANCELED"
	case ExecutionTypeRejected:
		return "REJECTED"
	case ExecutionTypeExpired:
		return "EXPIRED"
	}
	return "UNKNOWN"
}

// IsValid 判断 ExecutionType 是否为已定义的类型
func (e ExecutionType) IsValid() bool {
	switch e {
	case ExecutionTypeNew,
		ExecutionTypeTrade,
		ExecutionTypeCanceled,
		ExecutionTypeRejected,
		ExecutionTypeExpired:
		return true
	default:
		return false
	}
}

// ParseExecutionType 从字符串解析 ExecutionType (不区分大小写)
func ParseExecutionType(s string) (ExecutionType, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "NEW":
		return ExecutionTypeNew, nil
	case "TRADE":
		return ExecutionTypeTrade, nil
	case "CANCELED":
		return ExecutionTypeCanceled, nil
	case "REJECTED":
		return ExecutionTypeRejected, nil
	case "EXPIRED":
		return ExecutionTypeExpired, nil
	default:
		return ExecutionTypeUnknown, fmt.Errorf("unknown execution type: %s", s)
	}
}

const (
	// ExecutionTypeUnknown 未知
	ExecutionTypeUnknown ExecutionType = iota
	// ExecutionTypeNew 新订单
	ExecutionTypeNew
	// ExecutionTypeTrade 成交
	ExecutionTypeTrade
	// ExecutionTypeCanceled 已取消
	ExecutionTypeCanceled
	// ExecutionTypeRejected 已拒绝
	ExecutionTypeRejected
	// ExecutionTypeExpired 已过期
	ExecutionTypeExpired
)

// OrderState 订单状态: 1-OrderStateNew, 2-OrderStatePartiallyFilled, 3-OrderStateFilled, 4-OrderStateCanceled, 5-OrderStateRejected
type OrderState int

// String 返回字符串表示
func (o OrderState) String() string {
	switch o {
	case OrderStateNew:
		return "NEW"
	case OrderStatePartiallyFilled:
		return "PARTIALLY_FILLED"
	case OrderStateFilled:
		return "FILLED"
	case OrderStateCanceled:
		return "CANCELED"
	case OrderStateRejected:
		return "REJECTED"
	}
	return "UNKNOWN"
}

// IsValid 判断 OrderState 是否为已定义的类型
func (o OrderState) IsValid() bool {
	switch o {
	case OrderStateNew,
		OrderStatePartiallyFilled,
		OrderStateFilled,
		OrderStateCanceled,
		OrderStateRejected:
		return true
	default:
		return false
	}
}

// ParseOrderState 从字符串解析 OrderState (不区分大小写)
func ParseOrderState(s string) (OrderState, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "NEW":
		return OrderStateNew, nil
	case "PARTIALLY_FILLED":
		return OrderStatePartiallyFilled, nil
	case "FILLED":
		return OrderStateFilled, nil
	case "CANCELED":
		return OrderStateCanceled, nil
	case "REJECTED":
		return OrderStateRejected, nil
	default:
		return OrderStateUnknown, fmt.Errorf("unknown order state: %s", s)
	}
}

const (
	// OrderStateUnknown 未知
	OrderStateUnknown OrderState = iota
	// OrderStateNew 新订单
	OrderStateNew
	// OrderStatePartiallyFilled 部分成交
	OrderStatePartiallyFilled
	// OrderStateFilled 全部成交
	OrderStateFilled
	// OrderStateCanceled 已取消
	OrderStateCanceled
	// OrderStateRejected 已拒绝
	OrderStateRejected
)

// OrderType 订单类型: 1-OrderTypeMarket, 2-OrderTypeLimit
type OrderType int

// String 返回字符串表示
func (o OrderType) String() string {
	switch o {
	case OrderTypeMarket:
		return "MARKET"
	case OrderTypeLimit:
		return "LIMIT"
	}
	return "UNKNOWN"
}

// IsValid 判断 OrderType 是否为已定义的类型
func (o OrderType) IsValid() bool {
	switch o {
	case OrderTypeMarket,
		OrderTypeLimit:
		return true
	default:
		return false
	}
}

// ParseOrderType 从字符串解析 OrderType (不区分大小写)
func ParseOrderType(s string) (OrderType, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "MARKET":
		return OrderTypeMarket, nil
	case "LIMIT":
		return OrderTypeLimit, nil
	default:
		return OrderTypeUnknown, fmt.Errorf("unknown order type: %s", s)
	}
}

const (
	// OrderTypeUnknown 未知
	OrderTypeUnknown OrderType = iota
	// OrderTypeMarket 市价单
	OrderTypeMarket
	// OrderTypeLimit 限价单
	OrderTypeLimit
)

// PositionStatus 持仓状态:
// 1-NewPosition 新开仓, 2-OpeningPosition 开仓中
// 3-HoldingPosition 持仓中, 4-ClosingPosition 平仓中, 5-ClosedPosition 已平仓
type PositionStatus int

// String 返回字符串表示
func (p PositionStatus) String() string {
	switch p {
	case NewPosition:
		return "NEW"
	case OpeningPosition:
		return "OPENING"
	case HoldingPosition:
		return "HOLDING"
	case ClosingPosition:
		return "CLOSING"
	case ClosedPosition:
		return "CLOSED"
	}
	return "UNKNOWN"
}

// IsValid 判断 PositionStatus 是否为已定义的类型
func (p PositionStatus) IsValid() bool {
	switch p {
	case NewPosition,
		OpeningPosition,
		HoldingPosition,
		ClosingPosition,
		ClosedPosition:
		return true
	default:
		return false
	}
}

// ParsePositionStatus 从字符串解析 PositionStatus (不区分大小写)
func ParsePositionStatus(s string) (PositionStatus, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "NEW":
		return NewPosition, nil
	case "OPENING":
		return OpeningPosition, nil
	case "HOLDING":
		return HoldingPosition, nil
	case "CLOSING":
		return ClosingPosition, nil
	case "CLOSED":
		return ClosedPosition, nil
	default:
		return PositionStatusUnknown, fmt.Errorf("unknown position status: %s", s)
	}
}

const (
	// PositionStatusUnknown 未知
	PositionStatusUnknown PositionStatus = iota
	// NewPosition 新开仓
	NewPosition
	// OpeningPosition 开仓中
	OpeningPosition
	// HoldingPosition 持仓中
	HoldingPosition
	// ClosingPosition 平仓中
	ClosingPosition
	// ClosedPosition 已平仓
	ClosedPosition
)

// SideType 方向: 1-SideTypeBuy, 2-SideTypeSell
type SideType int

// String 返回字符串表示
func (s SideType) String() string {
	switch s {
	case SideTypeBuy:
		return "BUY"
	case SideTypeSell:
		return "SELL"
	}
	return "UNKNOWN"
}

// IsValid 判断 SideType 是否为已定义的类型
func (s SideType) IsValid() bool {
	switch s {
	case SideTypeBuy,
		SideTypeSell:
		return true
	default:
		return false
	}
}

// ParseSideType 从字符串解析 SideType (不区分大小写)
func ParseSideType(s string) (SideType, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "BUY":
		return SideTypeBuy, nil
	case "SELL":
		return SideTypeSell, nil
	default:
		return SideTypeUnknown, fmt.Errorf("unknown side type: %s", s)
	}
}

const (
	// SideTypeUnknown 未知
	SideTypeUnknown SideType = iota
	// SideTypeBuy 买入
	SideTypeBuy
	// SideTypeSell 卖出
	SideTypeSell
)

// PositionSide 持仓方向: 1-PositionSideLong, 2-PositionSideShort
type PositionSide int

// String 返回字符串表示
func (p PositionSide) String() string {
	switch p {
	case PositionSideLong:
		return "LONG"
	case PositionSideShort:
		return "SHORT"
	}
	return "UNKNOWN"
}

// IsValid 判断 PositionSide 是否为已定义的类型
func (p PositionSide) IsValid() bool {
	switch p {
	case PositionSideLong,
		PositionSideShort:
		return true
	default:
		return false
	}
}

// ParsePositionSide 从字符串解析 PositionSide (不区分大小写)
func ParsePositionSide(s string) (PositionSide, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "LONG":
		return PositionSideLong, nil
	case "SHORT":
		return PositionSideShort, nil
	default:
		return PositionSideUnknown, fmt.Errorf("unknown position side: %s", s)
	}
}

const (
	// PositionSideUnknown 未知
	PositionSideUnknown PositionSide = iota
	// PositionSideLong 多头
	PositionSideLong
	// PositionSideShort 空头
	PositionSideShort
)

// TimeInForce 成交条件: 1-TimeInForceGTC, 2-TimeInForceIOC, 3-TimeInForceFOK
type TimeInForce int

// String 返回字符串表示
func (t TimeInForce) String() string {
	switch t {
	case TimeInForceGTC:
		return "GTC"
	case TimeInForceIOC:
		return "IOC"
	case TimeInForceFOK:
		return "FOK"
	}
	return "UNKNOWN"
}

// IsValid 判断 TimeInForce 是否为已定义的类型
func (t TimeInForce) IsValid() bool {
	switch t {
	case TimeInForceGTC,
		TimeInForceIOC,
		TimeInForceFOK:
		return true
	default:
		return false
	}
}

// ParseTimeInForce 从字符串解析 TimeInForce (不区分大小写)
func ParseTimeInForce(s string) (TimeInForce, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "GTC":
		return TimeInForceGTC, nil
	case "IOC":
		return TimeInForceIOC, nil
	case "FOK":
		return TimeInForceFOK, nil
	default:
		return TimeInForceUnknown, fmt.Errorf("unknown time in force: %s", s)
	}
}

const (
	// TimeInForceUnknown 未知
	TimeInForceUnknown TimeInForce = iota
	// TimeInForceGTC 成交为止
	TimeInForceGTC
	// TimeInForceIOC 立即成交或取消
	TimeInForceIOC
	// TimeInForceFOK 全部成交或立即取消
	TimeInForceFOK
)

// OrderStatus 订单状态:
// 1-OrderStatusNew, 2-OrderStatusPartiallyFilled, 3-OrderStatusFilled,
// 4-OrderStatusCanceled, 5-OrderStatusPendingCancel, 6-OrderStatusRejected
type OrderStatus int

// String 返回字符串表示
func (o OrderStatus) String() string {
	switch o {
	case OrderStatusNew:
		return "NEW"
	case OrderStatusPartiallyFilled:
		return "PARTIALLY_FILLED"
	case OrderStatusFilled:
		return "FILLED"
	case OrderStatusCanceled:
		return "CANCELED"
	case OrderStatusPendingCancel:
		return "PENDING_CANCEL"
	case OrderStatusRejected:
		return "REJECTED"
	}
	return "UNKNOWN"
}

// IsValid 判断 OrderStatus 是否为已定义的类型
func (o OrderStatus) IsValid() bool {
	switch o {
	case OrderStatusNew,
		OrderStatusPartiallyFilled,
		OrderStatusFilled,
		OrderStatusCanceled,
		OrderStatusPendingCancel,
		OrderStatusRejected:
		return true
	default:
		return false
	}
}

// ParseOrderStatus 从字符串解析 OrderStatus (不区分大小写)
func ParseOrderStatus(s string) (OrderStatus, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "NEW":
		return OrderStatusNew, nil
	case "PARTIALLY_FILLED":
		return OrderStatusPartiallyFilled, nil
	case "FILLED":
		return OrderStatusFilled, nil
	case "CANCELED":
		return OrderStatusCanceled, nil
	case "PENDING_CANCEL":
		return OrderStatusPendingCancel, nil
	case "REJECTED":
		return OrderStatusRejected, nil
	default:
		return OrderStatusUnknown, fmt.Errorf("unknown order status: %s", s)
	}
}

const (
	// OrderStatusUnknown 未知
	OrderStatusUnknown OrderStatus = iota
	// OrderStatusNew 新订单
	OrderStatusNew
	// OrderStatusPartiallyFilled 部分成交
	OrderStatusPartiallyFilled
	// OrderStatusFilled 全部成交
	OrderStatusFilled
	// OrderStatusCanceled 已取消
	OrderStatusCanceled
	// OrderStatusPendingCancel 待取消
	OrderStatusPendingCancel
	// OrderStatusRejected 已拒绝
	OrderStatusRejected
)
