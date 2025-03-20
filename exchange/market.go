package exchange

import (
	"context"

	"github.com/go-gotop/gotop/types"
	"github.com/shopspring/decimal"
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

	// ConvertCoinToContract 将币转换为张(未处理标的物精度)
	// 参数：
	//   ctx: 上下文，用于控制请求超时、取消等。
	//   req: 请求参数对象，包含合约面值(ctVal)、数量(size)、市场类型(marketType)、价格(price)等信息。
	// 返回值：
	//   decimal.Decimal: 转换后的数量
	ConvertCoinToContract(ctx context.Context, req *ConvertSizeUnitRequest) (decimal.Decimal, error)

	// ConvertContractToCoin 将张转换为币(未处理标的物精度)
	// 参数：
	//   ctx: 上下文，用于控制请求超时、取消等。
	//   req: 请求参数对象，包含合约面值(ctVal)、数量(size)、市场类型(marketType)、价格(price)等信息。
	// 返回值：
	//   decimal.Decimal: 转换后的数量
	ConvertContractToCoin(ctx context.Context, req *ConvertSizeUnitRequest) (decimal.Decimal, error)

	// ConvertQuoteToContract 将报价转换为张(未处理标的物精度)
	// 参数：
	//   ctx: 上下文，用于控制请求超时、取消等。
	//   req: 请求参数对象，包含合约面值(ctVal)、数量(size)、市场类型(marketType)、价格(price)等信息。
	// 返回值：
	//   decimal.Decimal: 转换后的数量
	ConvertQuoteToContract(ctx context.Context, req *ConvertSizeUnitRequest) (decimal.Decimal, error)

	// ConvertContractToQuote 将张转换为报价(未处理标的物精度)
	// 参数：
	//   ctx: 上下文，用于控制请求超时、取消等。
	//   req: 请求参数对象，包含合约面值(ctVal)、数量(size)、市场类型(marketType)、价格(price)等信息。
	// 返回值：
	//   decimal.Decimal: 转换后的数量
	ConvertContractToQuote(ctx context.Context, req *ConvertSizeUnitRequest) (decimal.Decimal, error)
}

// ConvertSizeUnitRequest 转换数量单位请求参数
type ConvertSizeUnitRequest struct {
	// CtVal 合约面值
	CtVal decimal.Decimal
	// Size 数量
	Size decimal.Decimal
	// MarketType 市场类型
	MarketType types.MarketType
	// Price 价格
	Price decimal.Decimal
}

// GetDepthRequest 获取深度请求参数
type GetDepthRequest struct {
	// Symbol 交易对
	Symbol string
	// Level 深度级别
	Level int
	// Type 市场类型
	Type types.MarketType
}

// GetDepthResponse 获取深度响应
type GetDepthResponse struct {
	// Depth 深度
	Depth Depth
}

// Depth 市场深度
type Depth struct {
	// Asks 卖盘
	Asks []DepthItem
	// Bids 买盘
	Bids []DepthItem
}

// DepthItem 市场深度项
type DepthItem struct {
	// Price 价格
	Price decimal.Decimal
	// Amount 数量
	Amount decimal.Decimal
}

type GetMarkPriceKlineRequest struct {
}

type GetMarkPriceKlineResponse struct {
}

type Kline struct {
}
