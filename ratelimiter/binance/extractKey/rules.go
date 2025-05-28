package extractkey

import (
	"fmt"

	"github.com/go-gotop/gotop/ratelimiter"
	"github.com/go-gotop/gotop/types"
)

type TimesRuleKey struct {
}

func (r *TimesRuleKey) ExtractKeys(request ratelimiter.ExchangeRateLimiterRequest) string {
	return r.extractTimesRule(request)
}

// extractOrderRule 提取下单限流算法规则的键
// key = binance:{marketType}:{requestType}
func (r *TimesRuleKey) extractTimesRule(request ratelimiter.ExchangeRateLimiterRequest) string {
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

type WeightRuleKey struct {
}

func (r *WeightRuleKey) ExtractKeys(request ratelimiter.ExchangeRateLimiterRequest) string {
	return r.extractWeightRule(request)
}

// extractWeightRule 提取权重限流算法规则的键
// key = binance:{marketType}:request
func (r *WeightRuleKey) extractWeightRule(request ratelimiter.ExchangeRateLimiterRequest) string {
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
