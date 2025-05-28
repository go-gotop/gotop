package extractkey

import (
	"fmt"
	"os"

	"github.com/go-gotop/gotop/ratelimiter"
)

type TimesRedisKey struct {
}

func (r *TimesRedisKey) ExtractKeys(request ratelimiter.ExchangeRateLimiterRequest) string {
	return r.extractTimesRedisKey(request)
}

func (r *TimesRedisKey) extractTimesRedisKey(request ratelimiter.ExchangeRateLimiterRequest) string {
	switch request.RequestType {
	case ratelimiter.RequestTypeOrder,
		ratelimiter.RequestTypeCancelOrder:
		// 私有频道，根据userid限流
		return fmt.Sprintf("okx:%s:%s", request.RequestType, request.AccountID)
	case ratelimiter.RequestWsConnect:
		ip := request.IP
		if ip == "" {
			_ip := os.Getenv("HOST_IP")
			if _ip == "" {
				_ip = "unknown"
			}
			ip = _ip
		}
		// 公共频道，根据ip限流
		return fmt.Sprintf("okx:%s:%s", request.RequestType, ip)
	}
	return ""
}
