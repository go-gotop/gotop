package extractkey

import (
	"fmt"
	"os"

	"github.com/go-gotop/gotop/ratelimiter"
	"github.com/go-gotop/gotop/types"
)

type TimesRedisKey struct {
}

func (r *TimesRedisKey) ExtractKeys(request ratelimiter.ExchangeRateLimiterRequest) string {
	return r.extractTimesRedisKey(request)
}

// extractOrderRedisKey 提取下单限流算法规则的键
// key = binance:{marketType}:{requestType}:{accountID}
func (r *TimesRedisKey) extractTimesRedisKey(request ratelimiter.ExchangeRateLimiterRequest) string {
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
	return fmt.Sprintf("%s:%s:%s:%s", "binance", marketType, request.RequestType, request.AccountID)
}

type WeightRedisKey struct {
}

func (r *WeightRedisKey) ExtractKeys(request ratelimiter.ExchangeRateLimiterRequest) string {
	return r.extractWeightRedisKey(request)
}

// extractWeightRedisKey 提取权重限流算法规则的键
// key = binance:{marketType}:request:{ip}
func (r *WeightRedisKey) extractWeightRedisKey(request ratelimiter.ExchangeRateLimiterRequest) string {
	ip := request.IP
	if ip == "" {
		_ip := os.Getenv("HOST_IP")
		if _ip == "" {
			_ip = "unknown"
		}
		ip = _ip
	}

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

	return fmt.Sprintf("binance:%s:request:%s", marketType, ip)
}
