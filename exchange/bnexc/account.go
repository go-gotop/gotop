package bnexc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-gotop/gotop/exchange"
	"github.com/go-gotop/gotop/requests"
	bnexreq "github.com/go-gotop/gotop/requests/binance"
)

var _ exchange.AccountManager = &BnAccountManager{}

type BnAccountManager struct {
	client requests.RequestClient
}

func NewBnAccountManager() *BnAccountManager {
	adapter := bnexreq.NewBinanceAdapter()
	client := requests.NewClient()
	client.SetAdapter(adapter)
	return &BnAccountManager{
		client: client,
	}
}

func (b *BnAccountManager) GetBalances(ctx context.Context, authInfo exchange.AuthInfo) (*exchange.GetBalancesResponse, error) {
	apiUrl := BNEX_API_SPOT_URL + "/sapi/v1/capital/config/getall"

	resp, err := b.client.DoRequest(&requests.Request{
		Method: http.MethodGet,
		URL:    apiUrl,
		Auth: &requests.AuthInfo{
			APIKey:    authInfo.APIKey,
			SecretKey: authInfo.SecretKey,
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
		return nil, fmt.Errorf("get balances failed, status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var capitalInfo []bnCapitalRecoveryResponse
	err = json.Unmarshal(body, &capitalInfo)
	if err != nil {
		return nil, err
	}

	result := &exchange.GetBalancesResponse{
		Balances: make([]exchange.Balance, 0, len(capitalInfo)),
	}

	for _, asset := range capitalInfo {
		free, err := strconv.ParseFloat(asset.Free, 64)
		if err != nil {
			free = 0
		}

		locked, err := strconv.ParseFloat(asset.Locked, 64)
		if err != nil {
			locked = 0
		}

		// 只返回有余额的资产
		if free > 0 || locked > 0 {
			result.Balances = append(result.Balances, exchange.Balance{
				Asset:     asset.Coin,
				Available: free,
				Locked:    locked,
			})
		}
	}

	return result, nil
}

func (b *BnAccountManager) GetBalance(ctx context.Context, authInfo exchange.AuthInfo, asset string) (*exchange.GetBalanceResponse, error) {
	apiUrl := BNEX_API_SPOT_URL + "/sapi/v1/capital/config/getall"

	if asset == "" {
		return nil, errors.New("asset is required")
	}

	resp, err := b.client.DoRequest(&requests.Request{
		Method: http.MethodGet,
		URL:    apiUrl,
		Auth: &requests.AuthInfo{
			APIKey:    authInfo.APIKey,
			SecretKey: authInfo.SecretKey,
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
		return nil, fmt.Errorf("get balance failed, status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var capitalInfo []bnCapitalRecoveryResponse
	err = json.Unmarshal(body, &capitalInfo)
	if err != nil {
		return nil, err
	}

	for _, item := range capitalInfo {
		if item.Coin == asset {
			free, err := strconv.ParseFloat(item.Free, 64)
			if err != nil {
				free = 0
			}

			locked, err := strconv.ParseFloat(item.Locked, 64)
			if err != nil {
				locked = 0
			}

			return &exchange.GetBalanceResponse{
				Balance: exchange.Balance{
					Asset:     item.Coin,
					Available: free,
					Locked:    locked,
				},
			}, nil
		}
	}

	return nil, fmt.Errorf("asset %s not found", asset)
}
