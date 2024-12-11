package exchange

import (
	"context"

	"github.com/go-gotop/gotop/types"
	"github.com/shopspring/decimal"
)

// Exchange 交易所接口
type Exchange interface {
	// Name 交易所名称
	Name() string
	// Assets 交易所支持的资产
	Assets(ctx context.Context) ([]types.Asset, error)
	// 创建订单
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*CreateOrderResponse, error)
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	APIKey           string
	SecretKey        string
	// 秘钥 密码 (okex)
	Passphrase       string
	// 订单时间
	OrderTime        int64
	// 交易对
	Symbol           string
	// 客户端订单ID
	ClientOrderID    string
	// 订单方向
	Side             types.SideType
	// 订单类型
	OrderType        types.OrderType
	// 持仓方向
	PositionSide     types.PositionSide
	// 成交条件
	TimeInForce      types.TimeInForce
	// 市场类型
	MarketType       types.MarketType
	// 合约面值； 合约张数 = 合约数量 / 合约面值
	CtVal            decimal.Decimal
	// 订单数量
	Size             decimal.Decimal
	// 订单价格
	Price            decimal.Decimal
	// 是否是统一账户, 默认 false
	IsUnifiedAccount bool
}

type CreateOrderResponse struct {
	TransactTime     int64
	Symbol           string
	ClientOrderID    string
	OrderID          string
	Side             types.SideType
	State            types.OrderStatus
	PositionSide     types.PositionSide
	Price            decimal.Decimal
	OriginalQuantity decimal.Decimal
	ExecutedQuantity decimal.Decimal
}