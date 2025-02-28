package okxexc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-gotop/gotop/exchange"
	"github.com/go-gotop/gotop/requests"
	okxreq "github.com/go-gotop/gotop/requests/okx"
	"github.com/go-gotop/gotop/types"
	"github.com/shopspring/decimal"
)

var _ exchange.MarketDataProvider = &OkxMarketData{}

// OkxDepthResponse 获取市场深度响应
type OkxDepthResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
		Ts   string     `json:"ts"`
	} `json:"data"`
}

// OkxMarketData 提供市场行情数据相关的接口方法
type OkxMarketData struct {
	client requests.RequestClient
}

func NewOkxMarketData() *OkxMarketData {
	adapter := okxreq.NewOKXAdapter()
	client := requests.NewClient()
	client.SetAdapter(adapter)
	return &OkxMarketData{
		client: client,
	}
}

func (o *OkxMarketData) GetDepth(ctx context.Context, req *exchange.GetDepthRequest) (*exchange.GetDepthResponse, error) {
	apiUrl := OKX_API_BASE_URL + "/api/v5/market/books"
	if req.Type == types.MarketTypeFuturesUSDMargined {
		apiUrl = OKX_API_BASE_URL + "/api/v5/market/books-50"
	}
	resp, err := o.client.DoRequest(&requests.Request{
		Method: http.MethodGet,
		URL:    apiUrl,
		Params: map[string]any{
			"instId": req.Symbol,
			"sz":     fmt.Sprintf("%d", req.Level),
		},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var depthResp OkxDepthResponse
	err = json.Unmarshal(body, &depthResp)
	if err != nil {
		return nil, err
	}

	// 检查响应是否成功
	if len(depthResp.Data) == 0 {
		return nil, fmt.Errorf("failed to get depth data: %s", depthResp.Msg)
	}

	// 解析深度数据
	result := &exchange.GetDepthResponse{
		Depth: exchange.Depth{
			Asks: make([]exchange.DepthItem, 0, len(depthResp.Data[0].Asks)),
			Bids: make([]exchange.DepthItem, 0, len(depthResp.Data[0].Bids)),
		},
	}

	// 解析卖盘数据
	for _, ask := range depthResp.Data[0].Asks {
		if len(ask) >= 2 {
			price, err := decimal.NewFromString(ask[0])
			if err != nil {
				continue
			}
			amount, err := decimal.NewFromString(ask[1])
			if err != nil {
				continue
			}

			item := exchange.DepthItem{
				Price:  price,
				Amount: amount,
			}
			result.Depth.Asks = append(result.Depth.Asks, item)
		}
	}

	// 解析买盘数据
	for _, bid := range depthResp.Data[0].Bids {
		if len(bid) >= 2 {
			price, err := decimal.NewFromString(bid[0])
			if err != nil {
				continue
			}
			amount, err := decimal.NewFromString(bid[1])
			if err != nil {
				continue
			}

			item := exchange.DepthItem{
				Price:  price,
				Amount: amount,
			}
			result.Depth.Bids = append(result.Depth.Bids, item)
		}
	}

	return result, nil
}

func (o *OkxMarketData) GetMarkPriceKline(ctx context.Context, req *exchange.GetMarkPriceKlineRequest) (*exchange.GetMarkPriceKlineResponse, error) {
	// 暂未实现
	return nil, fmt.Errorf("方法未实现")
}
