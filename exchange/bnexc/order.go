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
	"github.com/shopspring/decimal"
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

	var err error
	size := req.Size

	if req.MarketType == types.MarketTypeFuturesCoinMargined || req.MarketType == types.MarketTypePerpetualCoinMargined {
		size, err = b.convertContractCoin("1", req.Symbol, req.Size)
		if err != nil {
			return nil, err
		}
	} else {
		// 下单的数量，统一进行向下取整
		size, err = b.sizePrecision(size, req.Symbol, "open")
		if err != nil {
			return nil, err
		}
	}

	params := map[string]any{
		"symbol":           req.Symbol.OriginalSymbol,
		"side":             req.Side.String(),
		"type":             req.OrderType.String(),
		"quantity":         size,
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
	return nil, nil
}

func (b *BnOrderManager) GetOrder(ctx context.Context, req *exchange.GetOrderRequest) (*exchange.GetOrderResponse, error) {
	return nil, nil
}

// 币转张(仅用于币本位合约) typ:1-币转张,2-张转币
func (b *BnOrderManager) convertContractCoin(typ string, symbol types.Symbol, size decimal.Decimal) (decimal.Decimal, error) {
	if symbol.Type != types.MarketTypeFuturesCoinMargined && symbol.Type != types.MarketTypePerpetualCoinMargined {
		return decimal.Zero, errors.New("invalid symbol market type")
	}

	if symbol.CtVal.Equal(decimal.Zero) {
		return decimal.Zero, errors.New("invalid symbol ct val")
	}

	price, err := b.getTickerPrice(symbol, symbol.Type)
	if err != nil {
		return decimal.Zero, err
	}

	switch typ {
	case "1":
		// 币转张：
		// 1. 获取TickerPrice
		// 2. 计算总币数usdt面值
		// 3. usdt面值除以合约面值
		// 4. 返回张数
		totalUsdt := size.Mul(price)
		num, err := b.sizePrecision(totalUsdt.Div(symbol.CtVal), symbol, "open")
		if err != nil {
			return decimal.Zero, err
		}
		return num, nil
	case "2":
		// 张转币：
		// 1. 获取TickerPrice
		// 2. 计算总张数usdt面值
		// 3. 总面值除以TickerPrice
		// 4. 返回币数
		totalUsdt := size.Mul(symbol.CtVal)
		num := totalUsdt.Div(price)
		return num, nil
	}
	return decimal.Zero, errors.New("invalid typ")
}

// 获取TickerPrice
func (b *BnOrderManager) getTickerPrice(symbol types.Symbol, marketType types.MarketType) (decimal.Decimal, error) {
	apiUrl := ""

	switch marketType {
	case types.MarketTypeFuturesCoinMargined, types.MarketTypePerpetualCoinMargined:
		apiUrl = BNEX_API_FUTURES_COIN_URL + "/dapi/v1/ticker/price"
	case types.MarketTypeFuturesUSDMargined, types.MarketTypePerpetualUSDMargined:
		apiUrl = BNEX_API_FUTURES_USD_URL + "/fapi/v1/ticker/price"
	case types.MarketTypeSpot:
		apiUrl = BNEX_API_SPOT_URL + "/api/v3/ticker/price"
	case types.MarketTypeMargin:
		apiUrl = BNEX_API_SPOT_URL + "/sapi/v1/margin/ticker/price"
	}

	var err error

	params := map[string]any{
		"symbol": symbol.OriginalSymbol,
	}

	resp, err := b.client.DoRequest(&requests.Request{
		Method: http.MethodGet,
		URL:    apiUrl,
		Params: params,
	})
	if err != nil {
		return decimal.Zero, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Zero, err
	}

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("get ticker price failed, status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var tickerPrice bnTickerPriceResponse
	err = json.Unmarshal(body, &tickerPrice)
	if err != nil {
		// 如果解析失败，尝试解析数组
		var tickerPriceArray []bnTickerPriceResponse
		err = json.Unmarshal(body, &tickerPriceArray)
		if err != nil {
			return decimal.Zero, err
		}
		tickerPrice = tickerPriceArray[0]
	}

	price, err := decimal.NewFromString(tickerPrice.Price)
	if err != nil {
		return decimal.Zero, err
	}
	return price, nil
}

// size 精度处理
// TODO: 没有做stepSize的限制
func (b *BnOrderManager) sizePrecision(size decimal.Decimal, symbol types.Symbol, opType string) (decimal.Decimal, error) {
	if symbol.MaxSize.Equal(decimal.Zero) {
		return decimal.Zero, errors.New("invalid symbol max size")
	}

	if symbol.MinSize.Equal(decimal.Zero) {
		return decimal.Zero, errors.New("invalid symbol min size")
	}

	orderQuantity := size
	if opType == "open" {
		// 向下取整到指定精度
		orderQuantity = orderQuantity.Truncate(symbol.SizePrecision)
	} else {
		// 四舍五入到指定精度
		orderQuantity = orderQuantity.Round(symbol.SizePrecision)
	}

	// 2. 限制最大值
	if orderQuantity.GreaterThan(symbol.MaxSize) {
		orderQuantity = symbol.MaxSize
	}

	// 3. 限制最小值
	if orderQuantity.LessThan(symbol.MinSize) {
		orderQuantity = symbol.MinSize
	}
	return orderQuantity, nil
}
