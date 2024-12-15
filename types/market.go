package types

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	// BinanceExchange 币安
	BinanceExchange = "BINANCE"
	// HuobiExchange 火币
	HuobiExchange = "HUOBI"
	// OkxExchange OKX
	OkxExchange = "OKX"
	// CoinBaseExchange CoinBase
	CoinBaseExchange = "COINBASE"
	// MockExchange 模拟
	MockExchange = "MOCK"

	// Maker 挂单
	Maker = "MAKER"
	// Taker 吃单
	Taker = "TAKER"

	// ByUser 用户
	ByUser = "USER"
	// BySystem 系统
	BySystem = "SYSTEM"
)

// MarketType 市场类型  
// 可扩展为：  
// 1. SPOT: 现货市场  
// 2. MARGIN: 杠杆/保证金市场  
// 3. FUTURES_USD_MARGINED: U本位期货（如 Binance USDT-M合约）
// 4. FUTURES_COIN_MARGINED: 币本位期货（如 Binance COIN-M合约）
// 5. PERPETUAL_USD_MARGINED: U本位永续合约
// 6. PERPETUAL_COIN_MARGINED: 币本位永续合约
// 7. OPTIONS: 期权市场
// 8. LEVERAGED_TOKENS: 杠杆代币
// 9. P2P: 点对点市场
// 10. ETF: ETF类产品市场
// 11. NFT: NFT数字藏品市场
type MarketType int

const (
	MarketTypeSpot MarketType = iota + 1
	MarketTypeMargin
	MarketTypeFuturesUSDMargined    // U本位期货
	MarketTypeFuturesCoinMargined   // 币本位期货
	MarketTypePerpetualUSDMargined  // U本位永续
	MarketTypePerpetualCoinMargined // 币本位永续
	MarketTypeOptions
	MarketTypeLeveragedTokens
	MarketTypeP2P
	MarketTypeETF
	MarketTypeNFT
)

// String 返回字符串表示
func (m MarketType) String() string {
	switch m {
	case MarketTypeSpot:
		return "SPOT"
	case MarketTypeMargin:
		return "MARGIN"
	case MarketTypeFuturesUSDMargined:
		return "FUTURES_USD_MARGINED"
	case MarketTypeFuturesCoinMargined:
		return "FUTURES_COIN_MARGINED"
	case MarketTypePerpetualUSDMargined:
		return "PERPETUAL_USD_MARGINED"
	case MarketTypePerpetualCoinMargined:
		return "PERPETUAL_COIN_MARGINED"
	case MarketTypeOptions:
		return "OPTIONS"
	case MarketTypeLeveragedTokens:
		return "LEVERAGED_TOKENS"
	case MarketTypeP2P:
		return "P2P"
	case MarketTypeETF:
		return "ETF"
	case MarketTypeNFT:
		return "NFT"
	default:
		return ""
	}
}

// IsValid 判断 MarketType 是否为已定义的类型
func (m MarketType) IsValid() bool {
	switch m {
	case MarketTypeSpot,
		MarketTypeMargin,
		MarketTypeFuturesUSDMargined,
		MarketTypeFuturesCoinMargined,
		MarketTypePerpetualUSDMargined,
		MarketTypePerpetualCoinMargined,
		MarketTypeOptions,
		MarketTypeLeveragedTokens,
		MarketTypeP2P,
		MarketTypeETF,
		MarketTypeNFT:
		return true
	default:
		return false
	}
}

// ParseMarketType 从字符串解析 MarketType (不区分大小写)
func ParseMarketType(s string) (MarketType, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "SPOT":
		return MarketTypeSpot, nil
	case "MARGIN":
		return MarketTypeMargin, nil
	case "FUTURES_USD_MARGINED":
		return MarketTypeFuturesUSDMargined, nil
	case "FUTURES_COIN_MARGINED":
		return MarketTypeFuturesCoinMargined, nil
	case "PERPETUAL_USD_MARGINED":
		return MarketTypePerpetualUSDMargined, nil
	case "PERPETUAL_COIN_MARGINED":
		return MarketTypePerpetualCoinMargined, nil
	case "OPTIONS":
		return MarketTypeOptions, nil
	case "LEVERAGED_TOKENS":
		return MarketTypeLeveragedTokens, nil
	case "P2P":
		return MarketTypeP2P, nil
	case "ETF":
		return MarketTypeETF, nil
	case "NFT":
		return MarketTypeNFT, nil
	default:
		return 0, fmt.Errorf("unknown market type: %s", s)
	}
}


// TradeEvent 成交事件
type TradeEvent struct {
	// ID 事件ID, 自增id由0开始
	// 如若需要用到该ID, 则需要自行实现
	ID uint64
	// Timestamp 事件时间
	Timestamp int64
	// Symbol 交易对
	Symbol string
	// Exchange 交易所
	Exchange string
	// Size 成交数量
	Size decimal.Decimal
	// Price 成交价格
	Price decimal.Decimal
	// Side 成交方向
	Side SideType
	// Type 市场类型
	Type MarketType
}

type Asset struct {
	// AssetName 资产名称
	AssetName string
	// Exchange 交易所
	Exchange string
	// MarketType 市场类型:
	// 1-MarketTypeSpot, 2-MarketTypeFutures, 3-MarketTypeOptions
	MarketType MarketType
	// Free 可用余额
	Free decimal.Decimal
	// Locked 锁定余额
	Locked decimal.Decimal
}

type Symbol struct {
	// OriginalSymbol 原标的物名称
	OriginalSymbol string
	// UnifiedSymbol 统一标的物名称
	UnifiedSymbol string
	// OriginalAsset 原资产名称
	OriginalAsset string
	// UnifiedAsset 统一资产名称
	UnifiedAsset string
	// 交易所
	Exchange string
	// Type 市场类型:
	// 1-MarketTypeSpot, 2-MarketTypeFutures, 3-MarketTypeOptions
	Type MarketType
	// Status 状态: ENABLED, DISABLED
	Status string
	// MinSize 最小头寸
	MinSize decimal.Decimal
	// MaxSize 最大头寸
	MaxSize decimal.Decimal
	// MinPrice 最小价格
	MinPrice decimal.Decimal
	// MaxPrice 最大价格
	MaxPrice decimal.Decimal
	// PricePrecision 价格精度
	PricePrecision int32
	// SizePrecision 头寸精度
	SizePrecision int32
	// CtVal 合约面值
	CtVal decimal.Decimal
	// CtMult 合约乘数
	CtMult decimal.Decimal
	// ListTime 上线时间
	ListTime int64
	// ExpTime 下线时间
	ExpTime int64
}