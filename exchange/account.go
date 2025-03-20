package exchange

import (
	"context"

	"github.com/shopspring/decimal"
)

// AccountManager 提供账户管理相关接口，如账户余额查询，账户信息获取
type AccountManager interface {
	// GetBalances 获取账户所有资产的余额信息
	// 参数：
	//   ctx: 上下文
	// 返回值：
	//   *GetBalancesResponse: 包含账户中各币种的可用余额和冻结(锁定)余额。
	//   error: 失败时返回错误信息。
	GetBalances(ctx context.Context, authInfo AuthInfo) (*GetBalancesResponse, error)

	// GetBalance 获取指定资产的余额信息
	// 参数：
	//   ctx: 上下文
	//   asset: 要查询的资产名称，如 "BTC"、"ETH"
	// 返回值：
	//   *GetBalanceResponse: 包含该资产可用余额、锁定余额信息
	//   error: 失败时返回错误信息。
	GetBalance(ctx context.Context, authInfo AuthInfo, asset string) (*GetBalanceResponse, error)
}

type AuthInfo struct {
	APIKey     string
	SecretKey  string
	Passphrase string
}

type GetBalancesResponse struct {
	Balances []Balance
}

type Balance struct {
	Asset     string
	Available decimal.Decimal
	Locked    decimal.Decimal
}

type GetBalanceResponse struct {
	Balance Balance
}
