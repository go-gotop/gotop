package okxexc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"

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

func (o *OkxMarketData) GetKline(ctx context.Context, req *exchange.GetKlineRequest) (*exchange.GetKlineResponse, error) {
	apiUrl := OKX_API_BASE_URL + "/api/v5/market/candles"

	params := map[string]any{
		"instId": req.Symbol,
		"bar":    req.Period,
	}

	if req.Start > 0 {
		params["before"] = fmt.Sprintf("%d", req.Start)
	}

	if req.End > 0 {
		params["after"] = fmt.Sprintf("%d", req.End)
	}

	if req.Limit > 0 {
		params["limit"] = fmt.Sprintf("%d", req.Limit)
	}

	resp, err := o.client.DoRequest(&requests.Request{
		Method: http.MethodGet,
		URL:    apiUrl,
		Params: params,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response okxKlineResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if response.Code != "0" {
		return nil, fmt.Errorf("operation failed, code: %s, message: %s", response.Code, response.Msg)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no data")
	}

	klines := response.Data

	result := &exchange.GetKlineResponse{
		Klines: make([]exchange.Kline, 0, len(klines)),
	}

	for _, kline := range klines {
		if len(kline) >= 7 {
			open, err := decimal.NewFromString(kline[1])
			if err != nil {
				return nil, err
			}
			high, err := decimal.NewFromString(kline[2])
			if err != nil {
				return nil, err
			}
			low, err := decimal.NewFromString(kline[3])
			if err != nil {
				return nil, err
			}
			close, err := decimal.NewFromString(kline[4])
			if err != nil {
				return nil, err
			}
			volume, err := decimal.NewFromString(kline[5])
			if err != nil {
				return nil, err
			}
			if req.MarketType == types.MarketTypeFuturesUSDMargined || req.MarketType == types.MarketTypePerpetualUSDMargined {
				volume, err = decimal.NewFromString(kline[6])
				if err != nil {
					return nil, err
				}
			}
			quoteVolume, err := decimal.NewFromString(kline[7])
			if err != nil {
				return nil, err
			}
			openTime, err := strconv.ParseInt(kline[0], 10, 64)
			if err != nil {
				return nil, err
			}
			confirm, err := strconv.ParseInt(kline[8], 10, 64)
			if err != nil {
				return nil, err
			}
			result.Klines = append(result.Klines, exchange.Kline{
				Symbol:      req.Symbol,
				Open:        open,
				High:        high,
				Low:         low,
				Close:       close,
				Volume:      volume,
				QuoteVolume: quoteVolume,
				OpenTime:    openTime,
				Confirm:     int(confirm),
			})
		}
	}

	// 正序
	sort.Slice(result.Klines, func(i, j int) bool {
		return result.Klines[i].OpenTime < result.Klines[j].OpenTime
	})

	return result, nil
}

func (o *OkxMarketData) GetMarkPriceKline(ctx context.Context, req *exchange.GetMarkPriceKlineRequest) (*exchange.GetMarkPriceKlineResponse, error) {
	// 暂未实现
	return nil, fmt.Errorf("方法未实现")
}

func ConvertCoinToContract(req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
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

func ConvertContractToCoin(req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
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

func ConvertQuoteToContract(req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
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

func ConvertContractToQuote(req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
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
