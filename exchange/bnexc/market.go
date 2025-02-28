package bnexc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	if req.Type == types.MarketTypeFuturesUSDMargined {
		apiUrl = BNEX_API_FUTURES_URL + "/fapi/v1/depth"
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

func (b *BnMarketData) GetMarkPriceKline(ctx context.Context, req *exchange.GetMarkPriceKlineRequest) (*exchange.GetMarkPriceKlineResponse, error) {
	return nil, nil
}
