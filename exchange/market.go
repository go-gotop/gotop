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

type GetDepthRequest struct {

}

type GetDepthResponse struct {

}

type Depth struct {

}

type GetMarkPriceKlineRequest struct {

}

type GetMarkPriceKlineResponse struct {

}

type Kline struct {

}
