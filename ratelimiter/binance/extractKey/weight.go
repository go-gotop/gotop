package extractkey

import (
	"fmt"

	"github.com/go-gotop/gotop/ratelimiter"
	"github.com/go-gotop/gotop/types"
)

type WeightKey struct {
}

func (w *WeightKey) ExtractKeys(request ratelimiter.ExchangeRateLimiterRequest) []string {
	keys := make([]string, 0)
	keys = append(keys, w.extractWeightKey(request))
	return keys
}

// extractWeightKey 提取权重限流算法规则的键
// key = binance:{marketType}:{requestType}:weight
func (w *WeightKey) extractWeightKey(request ratelimiter.ExchangeRateLimiterRequest) string {
	marketType := ""
	switch request.MarketType {
	case types.MarketTypeSpot, types.MarketTypeMargin:
		marketType = "spot"
	case types.MarketTypeFuturesUSDMargined,
		types.MarketTypeFuturesCoinMargined,
		types.MarketTypePerpetualUSDMargined,
		types.MarketTypePerpetualCoinMargined:
		marketType = "futures"
	default:
		return ""
	}
	return fmt.Sprintf("binance:%s:%s:weight", marketType, request.RequestType)
}
