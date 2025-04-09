package bnexc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-gotop/gotop/exchange"
	"github.com/go-gotop/gotop/requests"
	bnexreq "github.com/go-gotop/gotop/requests/binance"
	"github.com/go-gotop/gotop/types"
	"github.com/shopspring/decimal"
)

var _ exchange.MarketDataProvider = &BnMarketData{}

type BnDepthResponse struct {
	Ts  int64      `json:"T"`
	Bid [][]string `json:"bids"`
	Ask [][]string `json:"asks"`
}

type BnMarketData struct {
	client requests.RequestClient
}

func NewBnMarketData() *BnMarketData {
	adapter := bnexreq.NewBinanceAdapter()
	client := requests.NewClient()
	client.SetAdapter(adapter)
	return &BnMarketData{
		client: client,
	}
}

func (b *BnMarketData) GetDepth(ctx context.Context, req *exchange.GetDepthRequest) (*exchange.GetDepthResponse, error) {
	apiUrl := BNEX_API_SPOT_URL + "/api/v3/depth"
	if req.Type == types.MarketTypeFuturesUSDMargined || req.Type == types.MarketTypePerpetualUSDMargined {
		apiUrl = BNEX_API_FUTURES_USD_URL + "/fapi/v1/depth"
	} else if req.Type == types.MarketTypeFuturesCoinMargined || req.Type == types.MarketTypePerpetualCoinMargined {
		apiUrl = BNEX_API_FUTURES_COIN_URL + "/dapi/v1/depth"
	}
	resp, err := b.client.DoRequest(&requests.Request{
		Method: http.MethodGet,
		URL:    apiUrl,
		Params: map[string]any{
			"symbol": req.Symbol,
			"limit":  fmt.Sprintf("%d", req.Level),
		},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var depthResp BnDepthResponse
	err = json.Unmarshal(body, &depthResp)
	if err != nil {
		return nil, err
	}

	if len(depthResp.Ask) == 0 && len(depthResp.Bid) == 0 {
		return nil, fmt.Errorf("invalid symbol or empty depth data for: %s", req.Symbol)
	}

	result := &exchange.GetDepthResponse{
		Depth: exchange.Depth{
			Asks: make([]exchange.DepthItem, 0, len(depthResp.Ask)),
			Bids: make([]exchange.DepthItem, 0, len(depthResp.Bid)),
		},
	}

	for _, ask := range depthResp.Ask {
		if len(ask) >= 2 {
			price, err := decimal.NewFromString(ask[0])
			if err != nil {
				continue
			}
			amount, err := decimal.NewFromString(ask[1])
			if err != nil {
				continue
			}
			result.Depth.Asks = append(result.Depth.Asks, exchange.DepthItem{
				Price:  price,
				Amount: amount,
			})
		}
	}

	for _, bid := range depthResp.Bid {
		if len(bid) >= 2 {
			price, err := decimal.NewFromString(bid[0])
			if err != nil {
				continue
			}
			amount, err := decimal.NewFromString(bid[1])
			if err != nil {
				continue
			}
			result.Depth.Bids = append(result.Depth.Bids, exchange.DepthItem{
				Price:  price,
				Amount: amount,
			})
		}
	}

	return result, nil
}

func (b *BnMarketData) GetKline(ctx context.Context, req *exchange.GetKlineRequest) (*exchange.GetKlineResponse, error) {
	var apiUrl string

	if req.MarketType == types.MarketTypeSpot || req.MarketType == types.MarketTypeMargin {
		apiUrl = BNEX_API_SPOT_URL + "/api/v3/klines"
	} else if req.MarketType == types.MarketTypeFuturesUSDMargined || req.MarketType == types.MarketTypePerpetualUSDMargined {
		apiUrl = BNEX_API_FUTURES_USD_URL + "/fapi/v1/klines"
	} else if req.MarketType == types.MarketTypeFuturesCoinMargined || req.MarketType == types.MarketTypePerpetualCoinMargined {
		apiUrl = BNEX_API_FUTURES_COIN_URL + "/dapi/v1/klines"
	}

	params := map[string]any{
		"symbol":   req.Symbol,
		"interval": req.Period,
	}

	if req.Start != 0 {
		params["start"] = fmt.Sprintf("%d", req.Start)
	}

	if req.End != 0 {
		params["end"] = fmt.Sprintf("%d", req.End)
	}

	if req.Limit != 0 {
		params["limit"] = fmt.Sprintf("%d", req.Limit)
	}

	resp, err := b.client.DoRequest(&requests.Request{
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

	var klines [][]interface{}
	err = json.Unmarshal(body, &klines)
	if err != nil {
		return nil, err
	}

	result := &exchange.GetKlineResponse{
		Klines: make([]exchange.Kline, 0, len(klines)),
	}

	for _, kline := range klines {
		if len(kline) >= 8 {
			openTime := int64(kline[0].(float64))
			closeTime := int64(kline[6].(float64))
			open, err := decimal.NewFromString(kline[1].(string))
			if err != nil {
				return nil, err
			}
			high, err := decimal.NewFromString(kline[2].(string))
			if err != nil {
				return nil, err
			}
			low, err := decimal.NewFromString(kline[3].(string))
			if err != nil {
				return nil, err
			}
			close, err := decimal.NewFromString(kline[4].(string))
			if err != nil {
				return nil, err
			}
			volume, err := decimal.NewFromString(kline[5].(string))
			if err != nil {
				return nil, err
			}
			quoteVolume, err := decimal.NewFromString(kline[7].(string))
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
				CloseTime:   closeTime,
				Confirm:     1, // 默认设置为已确认
			})
		}
	}

	// 判断最后一条K线是否完结
	if len(result.Klines) > 0 {
		lastKline := result.Klines[len(result.Klines)-1]
		isComplete := isKlineComplete(lastKline.OpenTime, req.Period)
		result.IsLastComplete = isComplete

		// 如果K线未完结，设置Confirm为0
		if !isComplete {
			result.Klines[len(result.Klines)-1].Confirm = 0
		}
	}

	return result, nil
}

func (b *BnMarketData) GetMarkPriceKline(ctx context.Context, req *exchange.GetMarkPriceKlineRequest) (*exchange.GetMarkPriceKlineResponse, error) {
	return nil, nil
}

// 判断K线是否完结，通过比较开盘时间、收盘时间和周期
func isKlineComplete(openTime int64, period string) bool {
	var periodMillis int64

	// 解析周期字符串，转换为毫秒
	if len(period) >= 2 {
		unit := period[len(period)-1:]
		value, err := strconv.ParseInt(period[:len(period)-1], 10, 64)
		if err != nil {
			return false
		}

		switch unit {
		case "m":
			periodMillis = value * 60 * 1000
		case "h":
			periodMillis = value * 60 * 60 * 1000
		case "d":
			periodMillis = value * 24 * 60 * 60 * 1000
		case "w":
			periodMillis = value * 7 * 24 * 60 * 60 * 1000
		case "M":
			// 简化处理，月周期不支持，因为天数不固定
			return false
		default:
			return false
		}
	} else {
		return false
	}

	// 计算预期的收盘时间
	expectedCloseTime := openTime + periodMillis

	// 当前时间
	currentTime := time.Now().UnixMilli()

	// 如果当前时间超过了预期收盘时间，则认为K线已完结
	return currentTime >= expectedCloseTime
}

func ConvertCoinToContract(req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
	if req.MarketType != types.MarketTypeFuturesCoinMargined && req.MarketType != types.MarketTypePerpetualCoinMargined {
		return decimal.Zero, fmt.Errorf("invalid market type: %s", req.MarketType)
	}

	if req.Price.IsZero() {
		return decimal.Zero, fmt.Errorf("price is required")
	}

	if req.Size.IsZero() {
		return decimal.Zero, fmt.Errorf("size is required")
	}

	if req.CtVal.IsZero() {
		return decimal.Zero, fmt.Errorf("ctVal is required")
	}

	totalQuote := req.Size.Mul(req.Price)
	size := totalQuote.Div(req.CtVal)
	return size, nil
}

func ConvertContractToCoin(req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
	if req.MarketType != types.MarketTypeFuturesCoinMargined && req.MarketType != types.MarketTypePerpetualCoinMargined {
		return decimal.Zero, fmt.Errorf("invalid market type: %s", req.MarketType)
	}

	if req.Price.IsZero() {
		return decimal.Zero, fmt.Errorf("price is required")
	}

	if req.Size.IsZero() {
		return decimal.Zero, fmt.Errorf("size is required")
	}

	if req.CtVal.IsZero() {
		return decimal.Zero, fmt.Errorf("ctVal is required")
	}

	totalQuote := req.Size.Mul(req.CtVal)
	size := totalQuote.Div(req.Price)
	return size, nil
}

func ConvertQuoteToContract(req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
	if req.MarketType != types.MarketTypeFuturesUSDMargined && req.MarketType != types.MarketTypePerpetualUSDMargined {
		return decimal.Zero, fmt.Errorf("invalid market type: %s", req.MarketType)
	}

	if req.Size.IsZero() {
		return decimal.Zero, fmt.Errorf("size is required")
	}

	if req.CtVal.IsZero() {
		return decimal.Zero, fmt.Errorf("ctVal is required")
	}

	size := req.Size.Div(req.CtVal)
	return size, nil
}

func ConvertContractToQuote(req *exchange.ConvertSizeUnitRequest) (decimal.Decimal, error) {
	if req.MarketType != types.MarketTypeFuturesUSDMargined && req.MarketType != types.MarketTypePerpetualUSDMargined {
		return decimal.Zero, fmt.Errorf("invalid market type: %s", req.MarketType)
	}

	if req.Size.IsZero() {
		return decimal.Zero, fmt.Errorf("size is required")
	}

	if req.CtVal.IsZero() {
		return decimal.Zero, fmt.Errorf("ctVal is required")
	}

	totalQuote := req.Size.Mul(req.CtVal)
	return totalQuote, nil
}
