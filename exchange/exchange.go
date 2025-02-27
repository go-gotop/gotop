package exchange

// Exchange 交易所接口，整合了订单管理、市场数据和账户管理功能。
// 在实现时，若某些交易所不支持部分方法，可在运行时做特性检测或返回未实现的错误。
type Exchange interface {
	// Name 返回交易所的名称
	// 用于在日志、调试或多交易所管理时区分不同的交易所实例。
	Name() string

	OrderManager
	MarketDataProvider
	AccountManager
}
