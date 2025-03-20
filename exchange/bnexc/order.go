package bnexc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-gotop/gotop/exchange"
	"github.com/go-gotop/gotop/requests"
	bnexreq "github.com/go-gotop/gotop/requests/binance"
	"github.com/go-gotop/gotop/types"
)

var _ exchange.OrderManager = &BnOrderManager{}

type BnOrderManager struct {
	client requests.RequestClient
}

func NewBnOrderManager() *BnOrderManager {
	adapter := bnexreq.NewBinanceAdapter()
	client := requests.NewClient()
	client.SetAdapter(adapter)
	return &BnOrderManager{
		client: client,
	}
}

// CreateOrder 创建订单
func (b *BnOrderManager) CreateOrder(ctx context.Context, req *exchange.CreateOrderRequest) (*exchange.CreateOrderResponse, error) {
	// 下单数量单位过滤
	switch req.SizeUnit {
	case types.SizeUnitContract:
		if req.MarketType != types.MarketTypeFuturesCoinMargined && req.MarketType != types.MarketTypePerpetualCoinMargined {
			return nil, errors.New("create order error: unsupported contract size unit for coin margined futures")
		}
	case types.SizeUnitQuote:
		return nil, errors.New("create order error: unsupported quote size unit")
	case types.SizeUnitCoin:
		if req.MarketType == types.MarketTypeFuturesCoinMargined || req.MarketType == types.MarketTypePerpetualCoinMargined {
			return nil, errors.New("create order error: unsupported coin size unit for coin margined futures")
		}
	}

	apiUrl := ""

	switch req.MarketType {
	case types.MarketTypeFuturesUSDMargined, types.MarketTypePerpetualUSDMargined:
		apiUrl = BNEX_API_FUTURES_USD_URL + "/fapi/v1/order"
	case types.MarketTypeFuturesCoinMargined, types.MarketTypePerpetualCoinMargined:
		apiUrl = BNEX_API_FUTURES_COIN_URL + "/dapi/v1/order"
	case types.MarketTypeSpot:
		apiUrl = BNEX_API_SPOT_URL + "/api/v3/order"
	case types.MarketTypeMargin:
		apiUrl = BNEX_API_SPOT_URL + "/sapi/v1/margin/order"
	default:
		return nil, errors.New("invalid market type")
	}

	params := map[string]any{
		"symbol":           req.Symbol.OriginalSymbol,
		"side":             req.Side.String(),
		"type":             req.OrderType.String(),
		"quantity":         req.Size,
		"newClientOrderId": req.ClientOrderID,
		"newOrderRespType": "ACK",
	}

	// 杠杆开仓自动借贷，平仓不自动借贷
	if req.MarketType == types.MarketTypeMargin {
		if (req.Side == types.SideTypeBuy && req.PositionSide == types.PositionSideLong) ||
			(req.Side == types.SideTypeSell && req.PositionSide == types.PositionSideShort) {
			params["sideEffectType"] = "AUTO_BORROW_REPAY"
		}
	}

	if req.OrderType == types.OrderTypeLimit {
		params["timeInForce"] = "GTC"
		params["price"] = req.Price
	}

	if req.PositionSide != types.PositionSideUnknown {
		params["positionSide"] = req.PositionSide.String()
	}

	resp, err := b.client.DoRequest(&requests.Request{
		Method: http.MethodPost,
		URL:    apiUrl,
		Params: params,
		Auth: &requests.AuthInfo{
			APIKey:    req.APIKey,
			SecretKey: req.SecretKey,
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create order failed, status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var orderACK bnOrderACKResponse
	err = json.Unmarshal(body, &orderACK)
	if err != nil {
		return nil, err
	}

	return &exchange.CreateOrderResponse{
		OrderID:       fmt.Sprintf("%d", orderACK.OrderID),
		ClientOrderID: orderACK.ClientOrderID,
		Symbol:        orderACK.Symbol,
	}, nil
}

func (b *BnOrderManager) CancelOrder(ctx context.Context, req *exchange.CancelOrderRequest) (*exchange.CancelOrderResponse, error) {
	return nil, errors.New("not implemented")
}

func (b *BnOrderManager) GetOrder(ctx context.Context, req *exchange.GetOrderRequest) (*exchange.GetOrderResponse, error) {
	return nil, errors.New("not implemented")
}
