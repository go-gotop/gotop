package okx

import (
	"io"
	"net/http/httptest"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/go-gotop/gotop/requests"
)

func TestOKXAdapter_BuildRequest(t *testing.T) {
	adapter := &OKXAdapter{}

	tests := []struct {
		name        string
		request     *requests.Request
		wantMethod  string
		wantHeaders map[string]string
		wantBody    string
		wantErr     bool
	}{
		{
			name: "GET request with auth",
			request: &requests.Request{
				Method: http.MethodGet,
				URL:    "/api/v5/market/ticker",
				Params: map[string]interface{}{
					"instId": "BTC-USDT",
				},
				Auth: &requests.AuthInfo{
					APIKey:     "test-api-key",
					SecretKey:  "test-secret-key",
					Passphrase: "test-passphrase",
				},
			},
			wantMethod: http.MethodGet,
			wantHeaders: map[string]string{
				"Content-Type":          "application/json",
				"OK-ACCESS-KEY":         "test-api-key",
				"OK-ACCESS-PASSPHRASE":  "test-passphrase",
			},
			wantErr: false,
		},
		{
			name: "POST request with auth and params",
			request: &requests.Request{
				Method: http.MethodPost,
				URL:    "/api/v5/trade/order",
				Params: map[string]interface{}{
					"instId":  "BTC-USDT",
					"tdMode":  "cash",
					"side":    "buy",
					"ordType": "limit",
					"px":      "20000",
					"sz":      "0.01",
				},
				Auth: &requests.AuthInfo{
					APIKey:     "test-api-key",
					SecretKey:  "test-secret-key",
					Passphrase: "test-passphrase",
				},
			},
			wantMethod: http.MethodPost,
			wantHeaders: map[string]string{
				"Content-Type":          "application/json",
				"OK-ACCESS-KEY":         "test-api-key",
				"OK-ACCESS-PASSPHRASE":  "test-passphrase",
			},
			wantErr: false,
		},
		{
			name:    "Nil request",
			request: nil,
			wantErr: true,
		},
		{
			name: "Missing auth",
			request: &requests.Request{
				Method: http.MethodGet,
				URL:    "/api/v5/market/ticker",
			},
			wantErr: false,
			wantMethod: http.MethodGet,
			wantHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		},
		{
			name: "Empty method",
			request: &requests.Request{
				URL: "/api/v5/market/ticker",
					Auth: &requests.AuthInfo{
						APIKey:     "test-api-key",
						SecretKey:  "test-secret-key",
						Passphrase: "test-passphrase",
					},
			},
			wantErr: true,
		},
		{
			name: "Empty URL",
			request: &requests.Request{
				Method: http.MethodGet,
				Auth: &requests.AuthInfo{
					APIKey:     "test-api-key",
					SecretKey:  "test-secret-key",
					Passphrase: "test-passphrase",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepared, err := adapter.BuildRequest(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// 验证请求方法
			if prepared.Method != tt.wantMethod {
				t.Errorf("BuildRequest() method = %v, want %v", prepared.Method, tt.wantMethod)
			}

			// 验证必需的请求头
			for k, v := range tt.wantHeaders {
				if got := prepared.Headers.Get(k); got != v {
					t.Errorf("BuildRequest() header[%s] = %v, want %v", k, got, v)
				}
			}

			// 验证签名相关的请求头是否存在且格式正确
			if tt.request.Auth != nil && prepared.Headers.Get("OK-ACCESS-SIGN") == "" {
				t.Error("BuildRequest() missing OK-ACCESS-SIGN header")
			}
			if tt.request.Auth != nil && prepared.Headers.Get("OK-ACCESS-TIMESTAMP") == "" {
				t.Error("BuildRequest() missing OK-ACCESS-TIMESTAMP header")
			}

			// 验证时间戳格式
			timestamp := prepared.Headers.Get("OK-ACCESS-TIMESTAMP")
			if timestamp != "" {
				_, err = time.Parse("2006-01-02T15:04:05.000Z", timestamp)
				if err != nil {
					t.Errorf("Invalid timestamp format: %v", err)
				}
			}

			// 对于 POST 请求，验证请求体
			if tt.request != nil && tt.request.Method == http.MethodPost {
				if len(prepared.Body) == 0 {
					t.Error("BuildRequest() POST request should have body")
				}
				// 验证 JSON 格式
				var jsonBody map[string]interface{}
				if err := json.Unmarshal(prepared.Body, &jsonBody); err != nil {
					t.Errorf("BuildRequest() invalid JSON body: %v", err)
				}
				// 验证参数是否都在请求体中
				for k, v := range tt.request.Params {
					if jsonBody[k] != v {
						t.Errorf("BuildRequest() body parameter %s = %v, want %v", k, jsonBody[k], v)
					}
				}
			}
		})
	}
}

func TestSignMessage(t *testing.T) {
	tests := []struct {
		name      string
		signData  string
		secretKey string
		wantErr   bool
	}{
		{
			name:      "Valid signature",
			signData:  "2022-01-01T12:00:00.000ZGET/api/v5/market/ticker?instId=BTC-USDT",
			secretKey: "test-secret-key",
			wantErr:   false,
		},
		{
			name:      "Empty sign data",
			signData:  "",
			secretKey: "test-secret-key",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := signMessage(tt.signData, tt.secretKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("signMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Error("signMessage() returned empty signature")
			}
		})
	}
}

func TestOKXAdapter_Integration_NoAuth(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求头
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		// 验证请求方法和路径
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		expectedPath := "/api/v5/market/ticker"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// 验证查询参数
		if r.URL.Query().Get("instId") != "BTC-USDT" {
			t.Errorf("Expected instId=BTC-USDT, got %s", r.URL.Query().Get("instId"))
		}

		// 返回模拟响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": "0",
			"data": []map[string]interface{}{
				{
					"instId": "BTC-USDT",
					"last":   "50000",
				},
			},
		})
	}))
	defer server.Close()

	// 创建客户端和适配器
	client := requests.NewClient()
	adapter := NewOKXAdapter()
	client.SetAdapter(adapter)

	// 构建请求
	req := &requests.Request{
		Method: http.MethodGet,
		URL:    server.URL + "/api/v5/market/ticker",
		Params: map[string]interface{}{
			"instId": "BTC-USDT",
		},
	}

	// 发送请求
	resp, err := client.DoRequest(req)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}
	defer resp.Body.Close()

	// 验证响应
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["code"] != "0" {
		t.Errorf("Expected code 0, got %v", result["code"])
	}
}

func TestOKXAdapter_Integration_WithAuth(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证基本请求头
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		// 验证认证相关的请求头
		requiredHeaders := []string{
			"OK-ACCESS-KEY",
			"OK-ACCESS-SIGN",
			"OK-ACCESS-TIMESTAMP",
			"OK-ACCESS-PASSPHRASE",
		}
		for _, header := range requiredHeaders {
			if r.Header.Get(header) == "" {
				t.Errorf("Missing required header: %s", header)
			}
		}

		// 验证请求方法和路径
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		expectedPath := "/api/v5/trade/order"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// 验证请求体
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		var requestBody map[string]interface{}
		if err := json.Unmarshal(body, &requestBody); err != nil {
			t.Fatalf("Failed to unmarshal request body: %v", err)
		}

		expectedParams := map[string]interface{}{
			"instId":  "BTC-USDT",
			"tdMode":  "cash",
			"side":    "buy",
			"ordType": "limit",
			"px":      "50000",
			"sz":      "0.01",
		}
		for k, v := range expectedParams {
			if requestBody[k] != v {
				t.Errorf("Expected parameter %s to be %v, got %v", k, v, requestBody[k])
			}
		}

		// 返回模拟响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code": "0",
			"data": []map[string]interface{}{
				{
					"ordId":  "12345",
					"clOrdId": "abcdef",
					"tag":    "",
					"sCode":  "0",
					"sMsg":   "",
				},
			},
		})
	}))
	defer server.Close()

	// 创建客户端和适配器
	client := requests.NewClient()
	adapter := NewOKXAdapter()
	client.SetAdapter(adapter)

	// 构建请求
	req := &requests.Request{
		Method: http.MethodPost,
		URL:    server.URL + "/api/v5/trade/order",
		Params: map[string]interface{}{
			"instId":  "BTC-USDT",
			"tdMode":  "cash",
			"side":    "buy",
			"ordType": "limit",
			"px":      "50000",
			"sz":      "0.01",
		},
		Auth: &requests.AuthInfo{
			APIKey:     "test-api-key",
			SecretKey:  "test-secret-key",
			Passphrase: "test-passphrase",
		},
	}

	// 发送请求
	resp, err := client.DoRequest(req)
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}
	defer resp.Body.Close()

	// 验证响应
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["code"] != "0" {
		t.Errorf("Expected code 0, got %v", result["code"])
	}
} 