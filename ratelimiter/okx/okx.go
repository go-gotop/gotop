package okx

import (
	"context"
	"errors"

	"github.com/go-gotop/gotop/ratelimiter"
	"github.com/redis/go-redis/v9"
)

type okxRateLimiterManager struct {
	redisClient *redis.Client
}

func NewOkxRateLimiterManager(redisClient *redis.Client) ratelimiter.RateLimitManager[ratelimiter.ExchangeRateLimiterRequest] {
	return &okxRateLimiterManager{
		redisClient: redisClient,
	}
}

func (m *okxRateLimiterManager) PreCheck(ctx context.Context, request ratelimiter.ExchangeRateLimiterRequest) (ratelimiter.RateLimitDecision, error) {
	rateLimiters := make([]ratelimiter.RateLimiter[ratelimiter.ExchangeRateLimiterRequest], 0)

	switch request.RequestType {
	case ratelimiter.RequestTypeOrder,
		ratelimiter.RequestTypeNormal,
		ratelimiter.RequestWsConnect:
		rateLimiters = append(rateLimiters, NewGeneralRateLimiter(m.redisClient))
	default:
		return ratelimiter.RateLimitDecision{
			Allowed: false,
			Reason:  "unsupported request type",
		}, errors.New("unsupported request type")
	}

	for _, rateLimiter := range rateLimiters {
		decision, err := rateLimiter.Check(ctx, request)
		if err != nil {
			return ratelimiter.RateLimitDecision{
				Allowed: false,
				Reason:  err.Error(),
			}, err
		}
		if !decision.Allowed {
			return decision, nil
		}
	}

	return ratelimiter.RateLimitDecision{
		Allowed: true,
	}, nil
}
