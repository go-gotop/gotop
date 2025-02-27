package exchange

import (
	"context"
)

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

type CreateOrderRequest struct {

}

type CreateOrderResponse struct {

}



type CancelOrderRequest struct {

}

type CancelOrderResponse struct {

}

type GetOrderRequest struct {

}

type GetOrderResponse struct {

}

type Order struct {

}