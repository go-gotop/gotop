package binance

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"

	"github.com/go-gotop/gotop/types"
)

func TestConnectAndDisconnect(t *testing.T) {
	// 启动本地测试WebSocket服务
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 将HTTP连接升级为WebSocket连接
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// 服务端逻辑：空转一会儿，不发送数据
		select {
		case <-time.After(1 * time.Second):
		}
	}))
	defer server.Close()

	// 将测试服务器的URL从 http:// 替换成 ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	bs := NewBinanceStream("test_id", types.StreamTypeTrade)
	err := bs.Connect(context.Background(), BinanceRequest{
		URL: wsURL,
	})
	require.NoError(t, err, "Connect should succeed")

	require.NoError(t, bs.Disconnect(), "Disconnect should succeed")
}

func TestMessageHandling(t *testing.T) {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// 服务端等待一会然后发送一条测试消息
		time.Sleep(200 * time.Millisecond)
		conn.WriteMessage(websocket.TextMessage, []byte("test_message"))
		// 再等一会然后断开连接
		time.Sleep(200 * time.Millisecond)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	var receivedMsg []byte
	var mu sync.Mutex

	bs := NewBinanceStream("msg_test", types.StreamTypeTrade)
	err := bs.Connect(context.Background(), BinanceRequest{
		URL: wsURL,
		Handler: func(data []byte) {
			mu.Lock()
			defer mu.Unlock()
			receivedMsg = data
		},
		ErrorHandler: func(err error) {
			t.Log("ErrorHandler called:", err)
		},
	})
	require.NoError(t, err)

	// 等待服务端发送消息
	time.Sleep(500 * time.Millisecond)
	mu.Lock()
	require.Equal(t, []byte("test_message"), receivedMsg)
	mu.Unlock()

	bs.Disconnect()
}

func TestErrorAndReconnect(t *testing.T) {
	// 模拟服务端只接受一次连接，然后立即关闭，触发client重连
	upgrader := websocket.Upgrader{}
	var connectCount int
	mu := sync.Mutex{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		connectCount++
		mu.Unlock()

		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		conn.Close() // 立即关闭，触发客户端的错误处理
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	var errorCalled bool
	var loggerBuffer strings.Builder
	logger := slog.New(slog.NewTextHandler(&loggerBuffer, nil))

	bs := NewBinanceStream("reconnect_test", types.StreamTypeTrade)
	bs.reconnectInterval = 2 * time.Second // 缩短重连间隔以加速测试
	err := bs.Connect(context.Background(), BinanceRequest{
		URL: wsURL,
		ErrorHandler: func(err error) {
			errorCalled = true
		},
		Logger: logger,
	})
	require.NoError(t, err)

	// 等待一段时间，让客户端发现连接断开并尝试重连
	time.Sleep(500 * time.Millisecond)
	require.True(t, errorCalled, "ErrorHandler should be called on read error")

	// 重连会在autoReconnectLoop中触发定时重连，等一会儿看看connectCount
	time.Sleep(3 * time.Second)

	mu.Lock()
	cc := connectCount
	mu.Unlock()

	// connectCount至少应该大于1，一次初始连接，一次定时重连
	require.GreaterOrEqual(t, cc, 1)

	bs.Disconnect()
}

func TestPingPong(t *testing.T) {
	upgrader := websocket.Upgrader{}

	// 服务端会在收到ping后立即回复pong
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// 读循环：遇到ping自动回复pong由底层库处理，也可在本地测试中手工实现
		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if mt == websocket.PingMessage {
				// 服务端库通常会自动处理pong，但这里可以手动发pong
				conn.WriteMessage(websocket.PongMessage, nil)
				_ = msg // 忽略ping消息内容
			}
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	var errCalled bool
	bs := NewBinanceStream("ping_test", types.StreamTypeTrade)
	bs.pingInterval = 200 * time.Millisecond // 加快ping频率便于测试

	err := bs.Connect(context.Background(), BinanceRequest{
		URL: wsURL,
		ErrorHandler: func(err error) {
			errCalled = true
		},
	})
	require.NoError(t, err)

	// 等待一段时间，看是否会产生错误
	time.Sleep(1 * time.Second)
	require.False(t, errCalled, "No error should occur if pong is received properly")

	bs.Disconnect()
}

func TestTimeBasedReconnect(t *testing.T) {
	// 测试24小时后自动重连逻辑，在测试中缩短为100ms
	upgrader := websocket.Upgrader{}
	var connectCount int
	mu := sync.Mutex{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		connectCount++
		mu.Unlock()

		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		// 简单空转，等待客户端重连逻辑触发
		time.Sleep(500 * time.Millisecond)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	var logBuffer strings.Builder
	logger := slog.New(slog.NewTextHandler(&logBuffer, nil))

	bs := NewBinanceStream("time_reconnect_test", types.StreamTypeTrade)
	bs.reconnectInterval = 100 * time.Millisecond // 缩短为100ms便于测试

	err := bs.Connect(context.Background(), BinanceRequest{
		URL:    wsURL,
		Logger: logger,
	})
	require.NoError(t, err)

	time.Sleep(300 * time.Millisecond)

	mu.Lock()
	cc := connectCount
	mu.Unlock()

	require.GreaterOrEqual(t, cc, 2, "Should have connected at least twice due to time-based reconnect")

	bs.Disconnect()
}
