package extractkey

import (
	"fmt"

	"github.com/go-gotop/gotop/ratelimiter"
	"github.com/go-gotop/gotop/types"
)

type RuleKey struct {
}

func (r *RuleKey) ExtractKeys(request ratelimiter.ExchangeRateLimiterRequest) []string {
	keys := make([]string, 0)
	keys = append(keys, r.extractOrderRule(request))
	keys = append(keys, r.extractWeightRule(request))
	return keys
}

// extractOrderRule 提取下单限流算法规则的键
// key = binance:{marketType}:{requestType}
func (r *RuleKey) extractOrderRule(request ratelimiter.ExchangeRateLimiterRequest) string {
	marketType := ""
	switch request.MarketType {
	case types.MarketTypeSpot, types.MarketTypeMargin:
		marketType = "spot"
	case types.MarketTypeFuturesUSDMargined,
		types.MarketTypeFuturesCoinMargined,
		types.MarketTypePerpetualUSDMargined,
		types.MarketTypePerpetualCoinMargined:
		marketType = "futures"
	}
	return fmt.Sprintf("binance:%s:%s", marketType, request.RequestType)
}

// extractWeightRule 提取权重限流算法规则的键
// key = binance:{marketType}:request
func (r *RuleKey) extractWeightRule(request ratelimiter.ExchangeRateLimiterRequest) string {
	marketType := ""
	switch request.MarketType {
	case types.MarketTypeSpot, types.MarketTypeMargin:
		marketType = "spot"
	case types.MarketTypeFuturesUSDMargined,
		types.MarketTypeFuturesCoinMargined,
		types.MarketTypePerpetualUSDMargined,
		types.MarketTypePerpetualCoinMargined:
		marketType = "futures"
	}
	return fmt.Sprintf("binance:%s:request", marketType)
}
