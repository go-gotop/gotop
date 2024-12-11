package types

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
	return ""
}

const (
	// StrategyStatusRunning 运行中
	StrategyStatusRunning StrategyStatus = iota + 1
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
	return ""
}

const (
	// PriceDirectionUp 向上
	PriceDirectionUp PriceDirection = iota + 1
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
	return ""
}

const (
	// PositioningLevelSupport 支撑
	PositioningLevelSupport PositioningLevel = iota + 1
	// PositioningLevelResistance 阻力
	PositioningLevelResistance
)

