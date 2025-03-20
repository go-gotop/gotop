package okxexc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-gotop/gotop/exchange"
	"github.com/go-gotop/gotop/requests"
	okxreq "github.com/go-gotop/gotop/requests/okx"
	"github.com/go-gotop/gotop/types"
)

var _ exchange.OrderManager = &OkxOrderManager{}

type OkxOrderManager struct {
	client requests.RequestClient
}

func NewOkxOrderManager() *OkxOrderManager {
	adapter := okxreq.NewOKXAdapter()
	client := requests.NewClient()
	client.SetAdapter(adapter)
	return &OkxOrderManager{
		client: client,
	}
}

// CreateOrder 创建订单
func (o *OkxOrderManager) CreateOrder(ctx context.Context, req *exchange.CreateOrderRequest) (*exchange.CreateOrderResponse, error) {
	apiUrl := OKX_API_BASE_URL + "/api/v5/trade/order"

	switch req.SizeUnit {
	case types.SizeUnitContract:
		// 永续和交割的单位是张
		isPass := false
		if req.MarketType == types.MarketTypeFuturesUSDMargined || req.MarketType == types.MarketTypePerpetualUSDMargined || req.MarketType == types.MarketTypeFuturesCoinMargined || req.MarketType == types.MarketTypePerpetualCoinMargined {
			isPass = true
		}
		if !isPass {
			return nil, fmt.Errorf("create order error: unsupported contract size unit for %v", req.MarketType.String())
		}
	case types.SizeUnitCoin:
		isPass := false
		if req.MarketType == types.MarketTypeSpot {
			// 现货成交
			isPass = true
		} else if req.MarketType == types.MarketTypeMargin && req.OrderType == types.OrderTypeLimit {
			// 杠杆限价成交
			isPass = true
		} else if req.MarketType == types.MarketTypeMargin && req.Side == types.SideTypeSell && req.OrderType == types.OrderTypeMarket {
			// 杠杆市价卖出
			isPass = true
		}

		if !isPass {
			return nil, fmt.Errorf("create order error: unsupported coin size unit for %v %v %v", req.MarketType.String(), req.Side.String(), req.OrderType.String())
		}
	case types.SizeUnitQuote:
		isPass := false
		if req.MarketType == types.MarketTypeMargin && req.OrderType == types.OrderTypeMarket && req.Side == types.SideTypeBuy {
			isPass = true
		}
		if !isPass {
			return nil, fmt.Errorf("create order error: unsupported quote size unit for %v %v %v", req.MarketType.String(), req.Side.String(), req.OrderType.String())
		}
	default:
		return nil, fmt.Errorf("create order error: unsupported market type %v", req.MarketType.String())
	}

	params, err := o.toOrderParams(req)
	if err != nil {
		return nil, err
	}

	resp, err := o.client.DoRequest(&requests.Request{
		Method: http.MethodPost,
		URL:    apiUrl,
		Params: params,
		Auth: &requests.AuthInfo{
			APIKey:     req.APIKey,
			SecretKey:  req.SecretKey,
			Passphrase: req.Passphrase,
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
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var respData okxOrderResponse
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return nil, err
	}

	if respData.Code != "0" || len(respData.Data) == 0 || respData.Data[0].SCode != "0" {
		msg := respData.Msg
		code := respData.Code
		if len(respData.Data) > 0 {
			msg = respData.Data[0].SMsg
			code = respData.Data[0].SCode
		}
		return nil, fmt.Errorf("operation failed, code: %s, message: %s", code, msg)
	}

	return &exchange.CreateOrderResponse{
		OrderID:       respData.Data[0].OrdId,
		ClientOrderID: respData.Data[0].ClOrdId,
		Symbol:        respData.Data[0].SCode,
	}, nil
}

// CancelOrder 取消订单
func (o *OkxOrderManager) CancelOrder(ctx context.Context, req *exchange.CancelOrderRequest) (*exchange.CancelOrderResponse, error) {
	return nil, errors.New("not implemented")
}

// GetOrder 获取订单
func (o *OkxOrderManager) GetOrder(ctx context.Context, req *exchange.GetOrderRequest) (*exchange.GetOrderResponse, error) {
	return nil, errors.New("not implemented")
}

// 下单请求参数
func (o *OkxOrderManager) toOrderParams(req *exchange.CreateOrderRequest) (map[string]any, error) {
	params := map[string]any{
		"instId":  req.Symbol.OriginalSymbol,
		"clOrdId": req.ClientOrderID,
		"side":    toOkxSide(req.Side),
		"ordType": toOkxOrderType(req.OrderType),
		"sz":      req.Size.String(),
	}

	if req.MarketType == types.MarketTypeFuturesUSDMargined ||
		req.MarketType == types.MarketTypePerpetualUSDMargined ||
		req.MarketType == types.MarketTypeFuturesCoinMargined ||
		req.MarketType == types.MarketTypePerpetualCoinMargined {
		params["tdMode"] = toOkxPosMode(types.PosModeCross)
		params["posSide"] = toOkxPositionSide(req.PositionSide)
	} else if req.MarketType == types.MarketTypeSpot {
		params["tgtCcy"] = "base_ccy"
		params["tdMode"] = "cash"
	} else if req.MarketType == types.MarketTypeMargin {
		params["tdMode"] = toOkxPosMode(types.PosModeCross)
		params["ccy"] = "USDT"
	}

	// 限价单
	if req.OrderType == types.OrderTypeLimit && !req.Price.IsZero() {
		params["px"] = fmt.Sprintf("%v", req.Price)
	}

	return params, nil
}
