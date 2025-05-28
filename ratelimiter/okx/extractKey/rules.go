package extractkey

import (
	"fmt"

	"github.com/go-gotop/gotop/ratelimiter"
)

type TimesRuleKey struct {
}

func (r *TimesRuleKey) ExtractKeys(request ratelimiter.ExchangeRateLimiterRequest) string {
	return r.extractTimesRule(request)
}

func (r *TimesRuleKey) extractTimesRule(request ratelimiter.ExchangeRateLimiterRequest) string {
	return fmt.Sprintf("okx:%s", request.RequestType)
}
