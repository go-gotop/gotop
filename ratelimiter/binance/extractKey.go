package binance

import (
	"fmt"

	"github.com/go-gotop/gotop/ratelimiter"
	"github.com/go-gotop/gotop/types"
)

type Keys struct {
	RuleKey   string
	WeightKey string
	RedisKey  string
}

type KeyExtractor struct {
}

func NewKeyExtractor() *KeyExtractor {
	return &KeyExtractor{}
}

func (k *KeyExtractor) ExtractKeys(request ratelimiter.ExchangeRateLimiterRequest) []Keys {
	return []Keys{}
}

func (k *KeyExtractor) extractOrderKeys(request ratelimiter.ExchangeRateLimiterRequest) Keys {
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

	return Keys{
		RuleKey:   fmt.Sprintf("binance:%s:%s", marketType, request.RequestType),
		RedisKey:  fmt.Sprintf("binance:%s:%s:%s", marketType, request.RequestType, request.AccountID),
	}
}
func (k *KeyExtractor) extractWeightKeys(request ratelimiter.ExchangeRateLimiterRequest) Keys {
	return Keys{}
}
