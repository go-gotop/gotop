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

	var depthResp okxDepthResponse
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

func (o *OkxMarketData) ConvertCoinToContract(ctx context.Context, req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
	if req.CtVal.IsZero() {
		return decimal.Zero, fmt.Errorf("ctVal is required")
	}

	switch req.MarketType {
	case types.MarketTypeFuturesUSDMargined, types.MarketTypePerpetualUSDMargined:
		size := req.Size.Div(req.CtVal)
		return size, nil
	case types.MarketTypeFuturesCoinMargined, types.MarketTypePerpetualCoinMargined:
		if req.Price.IsZero() {
			return decimal.Zero, fmt.Errorf("price is required")
		}
		totalQuote := req.Size.Mul(req.Price)
		size := totalQuote.Div(req.CtVal)
		return size, nil
	default:
		return decimal.Zero, fmt.Errorf("invalid market type: %s", req.MarketType)
	}
}

func (o *OkxMarketData) ConvertContractToCoin(ctx context.Context, req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
	if req.CtVal.IsZero() {
		return decimal.Zero, fmt.Errorf("ctVal is required")
	}

	if req.Size.IsZero() {
		return decimal.Zero, fmt.Errorf("size is required")
	}

	switch req.MarketType {
	case types.MarketTypeFuturesUSDMargined, types.MarketTypePerpetualUSDMargined:
		size := req.Size.Mul(req.CtVal)
		return size, nil
	case types.MarketTypeFuturesCoinMargined, types.MarketTypePerpetualCoinMargined:
		if req.Price.IsZero() {
			return decimal.Zero, fmt.Errorf("price is required")
		}
		totalQuote := req.Size.Mul(req.CtVal)
		size := totalQuote.Div(req.Price)
		return size, nil
	default:
		return decimal.Zero, fmt.Errorf("invalid market type: %s", req.MarketType)
	}
}

func (o *OkxMarketData) ConvertQuoteToContract(ctx context.Context, req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
	if req.CtVal.IsZero() {
		return decimal.Zero, fmt.Errorf("ctVal is required")
	}

	if req.Size.IsZero() {
		return decimal.Zero, fmt.Errorf("size is required")
	}

	switch req.MarketType {
	case types.MarketTypeFuturesUSDMargined, types.MarketTypePerpetualUSDMargined:
		if req.Price.IsZero() {
			return decimal.Zero, fmt.Errorf("price is required")
		}
		size := req.Size.Div(req.Price).Div(req.CtVal)
		return size, nil
	case types.MarketTypeFuturesCoinMargined, types.MarketTypePerpetualCoinMargined:
		size := req.Size.Div(req.CtVal)
		return size, nil
	default:
		return decimal.Zero, fmt.Errorf("invalid market type: %s", req.MarketType)
	}
}

func (o *OkxMarketData) ConvertContractToQuote(ctx context.Context, req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
	if req.CtVal.IsZero() {
		return decimal.Zero, fmt.Errorf("ctVal is required")
	}

	if req.Size.IsZero() {
		return decimal.Zero, fmt.Errorf("size is required")
	}

	switch req.MarketType {
	case types.MarketTypeFuturesUSDMargined, types.MarketTypePerpetualUSDMargined:
		if req.Price.IsZero() {
			return decimal.Zero, fmt.Errorf("price is required")
		}
		totalQuote := req.Size.Mul(req.CtVal).Mul(req.Price)
		return totalQuote, nil
	case types.MarketTypeFuturesCoinMargined, types.MarketTypePerpetualCoinMargined:
		totalQuote := req.Size.Mul(req.CtVal)
		return totalQuote, nil
	default:
		return decimal.Zero, fmt.Errorf("invalid market type: %s", req.MarketType)
	}
}
