package exchange

import (
	"context"
)

// MarketDataProvider 提供市场行情数据相关的接口方法
type MarketDataProvider interface {
	// GetDepth 获取指定交易对的订单簿深度数据。
	// 参数：
	//   ctx: 上下文，用于控制请求超时、取消等。
	//   req: 请求参数对象，包含交易对(symbol)、请求深度(limit)等信息。
	// 返回值：
	//   *GetDepthResponse: 返回包含订单簿买卖盘数据的响应结构体。
	//   error: 如果请求过程中发生错误，返回错误信息。
	GetDepth(ctx context.Context, req *GetDepthRequest) (*GetDepthResponse, error)

	// GetMarkPriceKline 获取指定交易对的标记价格K线数据。
	// 一些合约交易所提供标记价格（Mark Price），该价格用于计算强平和资金费率。
	// 参数：
	//   ctx: 上下文，用于控制请求超时、取消等。
	//   req: 请求参数对象，包含交易对(symbol)、时间间隔(interval)、起止时间等信息。
	// 返回值：
	//   *GetMarkPriceKlineResponse: 包含K线数据点的列表。
	//   error: 请求失败时返回错误信息。
	GetMarkPriceKline(ctx context.Context, req *GetMarkPriceKlineRequest) (*GetMarkPriceKlineResponse, error)
}

// OrderManager 提供订单管理相关接口方法，如创建、撤销、查询订单等
type OrderManager interface {
	// CreateOrder 下订单
	// 参数：
	//   ctx: 上下文
	//   req: 包含下单所需的交易对(symbol)、方向(side)、价格(price)、数量(amount)、订单类型(type)等信息。
	// 返回值：
	//   *CreateOrderResponse: 包含订单ID、状态等返回信息。
	//   error: 失败时返回错误信息。
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error)

	// CancelOrder 撤销订单
	// 参数：
	//   ctx: 上下文
	//   req: 包含要撤销的订单ID或客户订单ID(ClientOrderID)等标识信息。
	// 返回值：
	//   *CancelOrderResponse: 包含撤单是否成功的状态信息。
	//   error: 失败时返回错误信息。
	CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CancelOrderResponse, error)

	// GetOrder 获取订单信息
	// 参数：
	//   ctx: 上下文
	//   req: 包含订单ID或其他查询所需字段
	// 返回值：
	//   *GetOrderResponse: 返回订单的详细信息，如状态、已成交数量、剩余数量等。
	//   error: 失败时返回错误信息。
	GetOrder(ctx context.Context, req *GetOrderRequest) (*GetOrderResponse, error)
}

// AccountManager 提供账户管理相关接口，如账户余额查询，账户信息获取
type AccountManager interface {
	// GetBalances 获取账户所有资产的余额信息
	// 参数：
	//   ctx: 上下文
	// 返回值：
	//   *GetBalancesResponse: 包含账户中各币种的可用余额和冻结(锁定)余额。
	//   error: 失败时返回错误信息。
	GetBalances(ctx context.Context) (*GetBalancesResponse, error)

	// GetBalance 获取指定资产的余额信息
	// 参数：
	//   ctx: 上下文
	//   asset: 要查询的资产名称，如 "BTC"、"ETH"
	// 返回值：
	//   *GetBalanceResponse: 包含该资产可用余额、锁定余额信息
	//   error: 失败时返回错误信息。
	GetBalance(ctx context.Context, asset string) (*GetBalanceResponse, error)
}

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
