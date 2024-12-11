package file

import (
	"os"
	"path/filepath"
	"context"
	"testing"
	"time"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/go-gotop/gotop/types"
)

func TestCSVFile_Subscribe(t *testing.T) {
	// 创建临时测试目录
	testDir := setupTestData(t)
	defer os.RemoveAll(testDir)

	// 使用固定的时间戳
	timestamp := int64(1657670400000)

	tests := []struct {
		ctx       context.Context
		name      string
		start     int64
		end       int64
		wantTicks int
		wantErr   bool
	}{
		{
			name:      "normal_case",
			ctx:       context.Background(),
			start:     timestamp,             // 使用精确的时间戳
			end:       timestamp + 1000,      // 包含两个数据点
			wantTicks: 2,
			wantErr:   false,
		},
		{
			name:      "invalid_time_range",
			ctx:       context.Background(),
			start:     timestamp + 2000,      // 在所有数据点之后
			end:       timestamp + 3000,
			wantTicks: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Test case: %s, start: %d, end: %d", tt.name, tt.start, tt.end)
			
			csvFile := NewCSVFile(testDir, tt.start, tt.end)
			require.NotNil(t, csvFile)

			var mu sync.Mutex
			tradeCount := 0
			trades := make([]types.TradeEvent, 0)
			var lastErr error
			done := make(chan struct{})
			var once sync.Once

			err := csvFile.Stream(
				tt.ctx,
				func(trade types.TradeEvent) {
					mu.Lock()
					defer mu.Unlock()
					t.Logf("Received tick: timestamp=%d, price=%s, size=%s, side=%v", 
						trade.Timestamp, trade.Price, trade.Size, trade.Side)
					trades = append(trades, trade)
					tradeCount++
					if tradeCount == tt.wantTicks {
						once.Do(func() {
							close(done)
						})
					}
				},
				func(err error) {
					t.Logf("Received error: %v", err)
					lastErr = err
					once.Do(func() {
						close(done)
					})
				},
			)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// 等待处理完成或超时
			select {
			case <-done:
				// 处理完成
			case <-time.After(time.Second):
				t.Fatal("timeout waiting for ticks")
			}

			mu.Lock()
			defer mu.Unlock()
			if tradeCount != tt.wantTicks {
				t.Errorf("Expected %d ticks, got %d. Ticks: %+v", tt.wantTicks, tradeCount, trades)
				for i, trade := range trades {
					t.Logf("Tick %d: timestamp=%d, price=%s, size=%s, side=%v", 
						i, trade.Timestamp, trade.Price, trade.Size, trade.Side)
				}
			}
			assert.Nil(t, lastErr)
		})
	}
}

func TestCSVFile_ProcessFile(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		wantTick *types.TradeEvent
		wantErr  bool
	}{
		{
			name: "valid_data",
			data: `trade_id,size,price,side,symbol,quote,traded_at
1,0.1,50000.00,BUY,BTCUSDT,USDT,1657670400000`,
			wantTick: &types.TradeEvent{
				Timestamp: 1657670400000,
				Price:     decimal.NewFromFloat(50000.00),
				Size:      decimal.NewFromFloat(0.1),
				Side:      types.SideTypeBuy,
			},
			wantErr: false,
		},
		{
			name: "invalid_price",
			data: `trade_id,size,price,side,symbol,quote,traded_at
1,0.1,invalid,BUY,BTCUSDT,USDT,1657670400000`,
			wantTick: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建临时文件
			tmpFile := createTempCSV(t, tt.data)
			defer os.Remove(tmpFile)

			csvFile := NewCSVFile("", 0, time.Now().UnixMilli())
			eventChan := make(chan types.TradeEvent, 1)

			err := csvFile.processFile(tmpFile, eventChan, 0, time.Now().UnixMilli())

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tt.wantTick != nil {
				select {
				case tick := <-eventChan:
					assert.Equal(t, tt.wantTick.Timestamp, tick.Timestamp)
					assert.True(t, tt.wantTick.Price.Equal(tick.Price))
					assert.True(t, tt.wantTick.Size.Equal(tick.Size))
					assert.Equal(t, tt.wantTick.Side, tick.Side)
				case <-time.After(time.Second):
					t.Fatal("timeout waiting for tick")
				}
			}
		})
	}
}

func TestParseSide(t *testing.T) {
	tests := []struct {
		name string
		side string
		want types.SideType
	}{
		{
			name: "buy",
			side: "BUY",
			want: types.SideTypeBuy,
		},
		{
			name: "sell",
			side: "SELL",
			want: types.SideTypeSell,
		},
		{
			name: "lowercase_buy",
			side: "buy",
			want: types.SideTypeBuy,
		},
		{
			name: "lowercase_sell",
			side: "sell",
			want: types.SideTypeSell,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSide(tt.side)
			assert.Equal(t, tt.want, got)
		})
	}
}

// 辅助函数
func setupTestData(t *testing.T) string {
	dir := t.TempDir()
	
	// 使用固定的时间戳
	data := `trade_id,size,price,side,symbol,quote,traded_at
1,0.1,50000.00,BUY,BTCUSDT,USDT,1657670400000
2,0.2,50100.00,SELL,BTCUSDT,USDT,1657670401000`

	// 直接在根目录创建文件
	filePath := filepath.Join(dir, "1657670400000.csv")
	err := os.WriteFile(filePath, []byte(data), 0644)
	require.NoError(t, err)

	t.Logf("Created test file at: %s with data:\n%s", filePath, data)
	
	// 验证文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Test file was not created: %v", err)
	}

	return dir
}

func createTempCSV(t *testing.T, data string) string {
	tmpFile, err := os.CreateTemp("", "test-*.csv")
	require.NoError(t, err)
	
	_, err = tmpFile.WriteString(data)
	require.NoError(t, err)
	
	err = tmpFile.Close()
	require.NoError(t, err)
	
	return tmpFile.Name()
} 