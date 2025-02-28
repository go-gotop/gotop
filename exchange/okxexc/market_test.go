package okxexc

import (
	"context"
	"testing"
	"time"

	"github.com/go-gotop/gotop/exchange"
	"github.com/go-gotop/gotop/types"
	"github.com/stretchr/testify/assert"
)

func TestOkxMarketData_GetDepth(t *testing.T) {
	// 跳过真实API测试，如需运行，请移除此行
	t.Skip("Skipping test that makes real API calls")

	tests := []struct {
		name        string
		symbol      string
		marketType  types.MarketType
		level       int
		expectError bool
	}{
		{
			name:        "Spot BTC-USDT Success",
			symbol:      "BTC-USDT",
			marketType:  types.MarketTypeSpot,
			level:       20,
			expectError: false,
		},
		{
			name:        "Futures BTC-USDT-SWAP Success",
			symbol:      "BTC-USDT-SWAP",
			marketType:  types.MarketTypeFuturesUSDMargined,
			level:       20,
			expectError: false,
		},
		{
			name:        "Invalid Symbol",
			symbol:      "INVALID-PAIR",
			marketType:  types.MarketTypeSpot,
			level:       20,
			expectError: true,
		},
	}

	marketData := NewOkxMarketData()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &exchange.GetDepthRequest{
				Symbol: tt.symbol,
				Type:   tt.marketType,
				Level:  tt.level,
			}

			resp, err := marketData.GetDepth(ctx, req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotNil(t, resp.Depth)

			// 验证深度数据
			assert.NotEmpty(t, resp.Depth.Asks, "Asks should not be empty")
			assert.NotEmpty(t, resp.Depth.Bids, "Bids should not be empty")

			// 验证价格和数量正确性
			if len(resp.Depth.Asks) > 0 {
				ask := resp.Depth.Asks[0]
				assert.False(t, ask.Price.IsZero(), "Ask price should not be zero")
				assert.False(t, ask.Amount.IsZero(), "Ask amount should not be zero")
			}

			if len(resp.Depth.Bids) > 0 {
				bid := resp.Depth.Bids[0]
				assert.False(t, bid.Price.IsZero(), "Bid price should not be zero")
				assert.False(t, bid.Amount.IsZero(), "Bid amount should not be zero")
			}

			// 验证买卖盘价格排序正确（卖盘价格应递增，买盘价格应递减）
			if len(resp.Depth.Asks) >= 2 {
				assert.True(t, resp.Depth.Asks[0].Price.LessThanOrEqual(resp.Depth.Asks[1].Price),
					"Ask prices should be in ascending order")
			}

			if len(resp.Depth.Bids) >= 2 {
				assert.True(t, resp.Depth.Bids[0].Price.GreaterThanOrEqual(resp.Depth.Bids[1].Price),
					"Bid prices should be in descending order")
			}
		})
	}
}

// TestOkxMarketData_GetDepth_Integration is a real integration test that can be run manually
func TestOkxMarketData_GetDepth_Integration(t *testing.T) {
	// 默认跳过，只在需要测试实际API时运行
	t.Skip("Integration test skipped by default")

	marketData := NewOkxMarketData()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &exchange.GetDepthRequest{
		Symbol: "BTC-USDT",
		Type:   types.MarketTypeSpot,
		Level:  20,
	}

	resp, err := marketData.GetDepth(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// 打印深度数据用于手动确认
	t.Logf("Received %d asks and %d bids", len(resp.Depth.Asks), len(resp.Depth.Bids))

	if len(resp.Depth.Asks) > 0 {
		t.Logf("First ask: price=%s, amount=%s", resp.Depth.Asks[0].Price.String(), resp.Depth.Asks[0].Amount.String())
	}

	if len(resp.Depth.Bids) > 0 {
		t.Logf("First bid: price=%s, amount=%s", resp.Depth.Bids[0].Price.String(), resp.Depth.Bids[0].Amount.String())
	}
}
