package requests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockAdapter 用于测试的模拟适配器
type mockAdapter struct {
	buildFunc func(*Request) (*PreparedRequest, error)
}

func (m *mockAdapter) BuildRequest(req *Request) (*PreparedRequest, error) {
	return m.buildFunc(req)
}

func TestClient_DoRequest(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	tests := []struct {
		name      string
		setupMock func() ExchangeAdapter
		request   *Request
		wantErr   bool
	}{
		{
			name: "Successful request",
			setupMock: func() ExchangeAdapter {
				return &mockAdapter{
					buildFunc: func(req *Request) (*PreparedRequest, error) {
						return &PreparedRequest{
							Method:  http.MethodGet,
							URL:     server.URL,
							Headers: http.Header{"Content-Type": []string{"application/json"}},
						}, nil
					},
				}
			},
			request: &Request{
				Method: http.MethodGet,
				URL:    "/test",
			},
			wantErr: false,
		},
		{
			name: "No adapter set",
			setupMock: func() ExchangeAdapter {
				return nil
			},
			request: &Request{
				Method: http.MethodGet,
				URL:    "/test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient()
			c.SetAdapter(tt.setupMock())

			resp, err := c.DoRequest(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Errorf("Failed to read response body: %v", err)
					return
				}

				if string(body) != `{"status":"ok"}` {
					t.Errorf("Unexpected response body: %s", string(body))
				}
			}
		})
	}
}

func TestClient_SetAdapter(t *testing.T) {
	c := NewClient()
	ma := &mockAdapter{}
	c.SetAdapter(ma)

	clientImpl := c.(*client)
	if clientImpl.adapter != ma {
		t.Error("SetAdapter did not set the adapter correctly")
	}
}

func TestClient_SetHTTPClient(t *testing.T) {
	c := NewClient()
	customClient := &http.Client{}
	c.SetHTTPClient(customClient)

	clientImpl := c.(*client)
	if clientImpl.httpClient != customClient {
		t.Error("SetHTTPClient did not set the HTTP client correctly")
	}
} 