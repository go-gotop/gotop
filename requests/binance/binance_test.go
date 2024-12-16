package binance

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/go-gotop/gotop/requests"
)

func TestBinanceAdapter_BuildRequest(t *testing.T) {
	adapter := NewBinanceAdapter()

	tests := []struct {
		name        string
		request     *requests.Request
		wantURL     string
		wantMethod  string
		wantHeaders map[string]string
		wantBody    string
		wantErr     bool
	}{
		{
			name: "GET request without auth",
			request: &requests.Request{
				Method: http.MethodGet,
				URL:    "/api/v3/ticker/price",
				Params: map[string]interface{}{
					"symbol": "BTCUSDT",
				},
			},
			wantURL:    "https://api.binance.com/api/v3/ticker/price?symbol=BTCUSDT",
			wantMethod: http.MethodGet,
			wantHeaders: map[string]string{
				"Accept":     "application/json",
				"User-Agent": "gotop/1.0",
			},
			wantBody: "",
			wantErr:  false,
		},
		{
			name: "POST request with auth",
			request: &requests.Request{
				Method: http.MethodPost,
				URL:    "/api/v3/order",
				Params: map[string]interface{}{
					"symbol": "BTCUSDT",
					"side":   "BUY",
					"type":   "LIMIT",
					"price":  "50000",
					"quantity": "0.001",
				},
				Auth: &requests.AuthInfo{
					APIKey:    "testApiKey",
					SecretKey: "testSecretKey",
				},
			},
			wantMethod: http.MethodPost,
			wantHeaders: map[string]string{
				"Accept":         "application/json",
				"User-Agent":     "gotop/1.0",
				"X-MBX-APIKEY":  "testApiKey",
				"Content-Type":   "application/x-www-form-urlencoded",
			},
			wantErr: false,
		},
		{
			name: "Invalid request - nil request",
			request: nil,
			wantErr: true,
		},
		{
			name: "Invalid request - empty URL",
			request: &requests.Request{
				Method: http.MethodGet,
				URL:    "",
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

			// 验证URL（对于带签名的请求，只验证基本部分）
			if tt.wantURL != "" {
				if prepared.URL != tt.wantURL {
					t.Errorf("BuildRequest() url = %v, want %v", prepared.URL, tt.wantURL)
				}
			}

			// 验证Headers
			for k, v := range tt.wantHeaders {
				if got := prepared.Headers.Get(k); got != v {
					t.Errorf("BuildRequest() header[%s] = %v, want %v", k, got, v)
				}
			}

			// 对于签名请求的特殊验证
			if tt.request.Auth != nil {
				// 验证是否包含timestamp和signature
				u, err := url.Parse(prepared.URL)
				if err != nil {
					t.Errorf("Failed to parse URL: %v", err)
					return
				}
				
				var params url.Values
				if tt.request.Method == http.MethodPost {
					params, err = url.ParseQuery(string(prepared.Body))
					if err != nil {
						t.Errorf("Failed to parse body: %v", err)
						return
					}
				} else {
					params = u.Query()
				}

				if params.Get("timestamp") == "" {
					t.Error("Missing timestamp in signed request")
				}
				if params.Get("signature") == "" {
					t.Error("Missing signature in signed request")
				}
			}
		})
	}
}

func TestBinanceAdapter_sign(t *testing.T) {
	tests := []struct {
		name    string
		message string
		secret  string
		want    string
	}{
		{
			name:    "Basic signature",
			message: "symbol=BTCUSDT&timestamp=1578963600000",
			secret:  "NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j",
			want:    "ea74d1a0ec7fff67761ee4dab78e249e73d84887f062253e4c83437a56a33778",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sign(tt.message, tt.secret)
			if got != tt.want {
				t.Errorf("sign() = %v, want %v", got, tt.want)
			}
		})
	}
} 