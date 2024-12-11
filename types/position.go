package types

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
	return ""
}

const (
	// OrderTypeMarket 市价单
	OrderTypeMarket OrderType = iota + 1
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
	return ""
}

const (
	// NewPosition 新开仓
	NewPosition PositionStatus = iota + 1
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
	return ""
}

const (
	// SideTypeBuy 买入
	SideTypeBuy SideType = iota + 1
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
	return ""
}

const (
	// PositionSideLong 多头
	PositionSideLong PositionSide = iota + 1
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
	return ""
}

const (
	// TimeInForceGTC 成交为止
	TimeInForceGTC TimeInForce = iota + 1
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
	}
	return ""
}

const (
	// OrderStatusNew 新订单
	OrderStatusNew OrderStatus = iota + 1
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
