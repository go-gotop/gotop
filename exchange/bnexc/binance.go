package bnexc

const (
	// 现货&杠杆
	BNEX_API_SPOT_URL = "https://api.binance.com"
	// U本位合约
	BNEX_API_FUTURES_USD_URL = "https://fapi.binance.com"
	// 币本位合约
	BNEX_API_FUTURES_COIN_URL = "https://dapi.binance.com"
)

// bnOrderACKResponse 币安下单返回响应(下单最快返回)
type bnOrderACKResponse struct {
	// 用户自定义的订单号
	ClientOrderID string `json:"clientOrderId"`
	// 系统订单号
	OrderID int64 `json:"orderId"`
	// 交易对
	Symbol string `json:"symbol"`
}

// bnSpotAccountInfoResponse 币安现货账户信息响应
type bnSpotAccountInfoResponse struct {
	// 手续费率
	MakerCommission int64 `json:"makerCommission"`
	// 手续费率
	TakerCommission int64 `json:"takerCommission"`
	// 手续费率
	BuyerCommission int64 `json:"buyerCommission"`
	// 手续费率
	SellerCommission int64 `json:"sellerCommission"`
	// 手续费率
	CommissionRates struct {
		// 手续费率
		Maker string `json:"maker"`
		// 手续费率
		Taker string `json:"taker"`
		// 手续费率
		Buyer string `json:"buyer"`
		// 手续费率
		Seller string `json:"seller"`
	} `json:"commissionRates"`
	// 是否可以交易
	CanTrade bool `json:"canTrade"`
	// 是否可以提现
	CanWithdraw bool `json:"canWithdraw"`
	// 是否可以充值
	CanDeposit bool `json:"canDeposit"`
	// 是否是中介
	Brokered bool `json:"brokered"`
	// 是否需要自定义交易预防
	RequireSelfTradePrevention bool `json:"requireSelfTradePrevention"`
	// 是否需要预防市价单
	PreventSor bool `json:"preventSor"`
	// 更新时间
	UpdateTime int64 `json:"updateTime"`
	// 账户类型
	AccountType string `json:"accountType"`
	// 余额
	Balances []struct {
		// 资产
		Asset string `json:"asset"`
		// 可用余额
		Free string `json:"free"`
		// 锁定余额
		Locked string `json:"locked"`
	} `json:"balances"`
	// 权限
	Permissions []string `json:"permissions"`
	// 用户ID
	Uid int64 `json:"uid"`
}

// bnMarginAccountInfoResponse 币安杠杆账户信息响应
type bnMarginAccountInfoResponse struct {
	// 是否已开户
	Created bool `json:"created"`
	// 是否可以借币
	BorrowEnabled bool `json:"borrowEnabled"`
	// 杠杆率
	MarginLevel string `json:"marginLevel"`
	// 抵押率
	CollateralMarginLevel string `json:"collateralMarginLevel"`
	// 总资产
	TotalAssetOfBtc string `json:"totalAssetOfBtc"`
	// 总负债
	TotalLiabilityOfBtc string `json:"totalLiabilityOfBtc"`
	// 总净资产
	TotalNetAssetOfBtc string `json:"totalNetAssetOfBtc"`
	// 总抵押资产
	TotalCollateralValueInUSDT string `json:"totalCollateralValueInUSDT"`
	// 总未实现亏损
	TotalOpenOrderLossInUSDT string `json:"totalOpenOrderLossInUSDT"`
	// 是否可以交易
	TradeEnabled bool `json:"tradeEnabled"`
	// 是否可以转入
	TransferInEnabled bool `json:"transferInEnabled"`
	// 是否可以转出
	TransferOutEnabled bool `json:"transferOutEnabled"`
	// 账户类型
	AccountType string `json:"accountType"`
	// 用户资产
	UserAssets []struct {
		// 资产
		Asset string `json:"asset"`
		// 借币
		Borrowed string `json:"borrowed"`
		// 可用余额
		Free string `json:"free"`
		// 利息
		Interest string `json:"interest"`
		// 锁定余额
		Locked string `json:"locked"`
		// 净资产
		NetAsset string `json:"netAsset"`
	} `json:"userAssets"`
}

// bnFuturesAccountInfoResponse 币安U本位合约账户信息响应
type bnFuturesAccountInfoResponse struct {
	// 当前所需起始保证金总额(存在逐仓请忽略), 仅计算usdt资产
	TotalInitialMargin string `json:"totalInitialMargin"`
	// 维持保证金总额, 仅计算usdt资产
	TotalMaintMargin string `json:"totalMaintMargin"`
	// 持仓未实现盈亏总额, 仅计算usdt资产
	TotalUnrealizedProfit string `json:"totalUnrealizedProfit"`
	// 保证金总余额, 仅计算usdt资产
	TotalMarginBalance string `json:"totalMarginBalance"`
	// 持仓所需起始保证金(基于最新标记价格), 仅计算usdt资产
	TotalPositionInitialMargin string `json:"totalPositionInitialMargin"`
	// 当前挂单所需起始保证金(基于最新标记价格), 仅计算usdt资产
	TotalOpenOrderInitialMargin string `json:"totalOpenOrderInitialMargin"`
	// 全仓账户余额, 仅计算usdt资产
	TotalCrossWalletBalance string `json:"totalCrossWalletBalance"`
	// 全仓持仓未实现盈亏总额, 仅计算usdt资产
	TotalCrossUnPnl string `json:"totalCrossUnPnl"`
	// 可用余额, 仅计算usdt资产
	AvailableBalance string `json:"availableBalance"`
	// 最大可转出余额, 仅计算usdt资产
	MaxWithdrawAmount string `json:"maxWithdrawAmount"`
	// 资产
	Assets []struct {
		// 资产
		Asset string `json:"asset"`
		// 余额
		WalletBalance string `json:"walletBalance"`
		// 未实现盈亏
		UnrealizedProfit string `json:"unrealizedProfit"`
		// 保证金余额
		MarginBalance string `json:"marginBalance"`
		// 维持保证金
		MaintMargin string `json:"maintMargin"`
		// 当前所需起始保证金
		InitialMargin string `json:"initialMargin"`
		// 持仓所需起始保证金(基于最新标记价格)
		PositionInitialMargin string `json:"positionInitialMargin"`
		// 全仓账户余额
		CrossWalletBalance string `json:"crossWalletBalance"`
		// 全仓持仓未实现盈亏
		CrossUnPnl string `json:"crossUnPnl"`
		// 可用余额
		AvailableBalance string `json:"availableBalance"`
		// 最大可转出余额
		MaxWithdrawAmount string `json:"maxWithdrawAmount"`
		// 更新时间
		UpdateTime int64 `json:"updateTime"`
	} `json:"assets"`
	// 持仓
	Positions []struct {
		// 资产
		Asset string `json:"asset"`
		// 余额
		WalletBalance string `json:"walletBalance"`
		// 持仓未实现盈亏
		UnrealizedProfit string `json:"unrealizedProfit"`
		// 保证金余额
		MarginBalance string `json:"marginBalance"`
		// 维持保证金
		MaintMargin string `json:"maintMargin"`
		// 当前所需起始保证金
		InitialMargin string `json:"initialMargin"`
		// 持仓所需起始保证金(基于最新标记价格)
		PositionInitialMargin string `json:"positionInitialMargin"`
		// 挂单所需起始保证金
		OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
		// 全仓账户余额
		CrossWalletBalance string `json:"crossWalletBalance"`
		// 全仓持仓未实现盈亏
		CrossUnPnl string `json:"crossUnPnl"`
		// 可用余额
		AvailableBalance string `json:"availableBalance"`
		// 最大可转出余额
		MaxWithdrawAmount string `json:"maxWithdrawAmount"`
		// 更新时间
		UpdateTime int64 `json:"updateTime"`
	} `json:"positions"`
}

// bnFuturesCoinAccountInfoResponse 币安币本位合约账户信息响应
type bnFuturesCoinAccountInfoResponse struct {
	// 资产内容
	Assets []struct {
		// 资产名
		Asset string `json:"asset"`
		// 账户余额
		WalletBalance string `json:"walletBalance"`
		// 全部持仓未实现盈亏
		UnrealizedProfit string `json:"unrealizedProfit"`
		// 维持保证金
		MaintMargin string `json:"maintMargin"`
		// 当前所需起始保证金(按最新标标记价格)
		InitialMargin string `json:"initialMargin"`
		// 当前所需持仓起始保证金(按最新标标记价格)
		PositionInitialMargin string `json:"positionInitialMargin"`
		// 当前所需挂单起始保证金(按最新标标记价格)
		OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
		// 最大可提款金额
		MaxWithdrawAmount string `json:"maxWithdrawAmount"`
		// 可用于全仓的账户余额
		CrossWalletBalance string `json:"crossWalletBalance"`
		// 所有全仓持仓的未实现盈亏
		CrossUnPnl string `json:"crossUnPnl"`
		// 可用下单余额
		AvailableBalance string `json:"availableBalance"`
		// 更新时间
		UpdateTime int64 `json:"updateTime"`
	} `json:"assets"`
	// 头寸
	Positions []struct {
		// 交易对
		Symbol string `json:"symbol"`
		// 持仓数量
		PositionAmt string `json:"positionAmt"`
		// 当前所需起始保证金(按最新标标记价格)
		InitialMargin string `json:"initialMargin"`
		// 持仓维持保证金
		MaintMargin string `json:"maintMargin"`
		// 持仓未实现盈亏
		UnrealizedProfit string `json:"unrealizedProfit"`
		// 当前所需持仓起始保证金(按最新标标记价格)
		PositionInitialMargin string `json:"positionInitialMargin"`
		// 当前所需挂单起始保证金(按最新标标记价格)
		OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
		// 杠杆倍率
		Leverage string `json:"leverage"`
		// 是否是逐仓模式
		Isolated bool `json:"isolated"`
		// 持仓方向
		PositionSide string `json:"positionSide"`
		// 平均持仓成本
		EntryPrice string `json:"entryPrice"`
		// 盈亏平衡价
		BreakEvenPrice string `json:"breakEvenPrice"`
		// 当前杠杆下最大可开仓数(标的数量)
		MaxQty string `json:"maxQty"`
		// 最新更新时间
		UpdateTime int64 `json:"updateTime"`
	} `json:"positions"`
	// 是否可以入金
	CanDeposit bool `json:"canDeposit"`
	// 是否可以交易
	CanTrade bool `json:"canTrade"`
	// 是否可以出金
	CanWithdraw bool `json:"canWithdraw"`
	// 手续费等级
	FeeTier int `json:"feeTier"`
	// 更新时间
	UpdateTime int64 `json:"updateTime"`
}

// bnCapitalRecoveryResponse 币安资金归集响应
type bnCapitalRecoveryResponse struct {
	// 资产名
	Coin string `json:"coin"`
	// 是否可以充值
	DepositAllEnable bool `json:"depositAllEnable"`
	// 可用余额
	Free string `json:"free"`
	// 冻结余额
	Freeze string `json:"freeze"`
	// 可申购余额
	Ipoable string `json:"ipoable"`
	// 可申购余额
	Ipoing string `json:"ipoing"`
	// 是否是合法货币
	IsLegalMoney bool `json:"isLegalMoney"`
	// 锁定余额
	Locked string `json:"locked"`
	// 资产名
	Name string `json:"name"`
	// 网络列表
	NetworkList []struct {
		// 地址正则
		AddressRegex string `json:"addressRegex"`
		// 资产名
		Coin string `json:"coin"`
		// 充值描述
		DepositDesc string `json:"depositDesc"`
		// 是否可以充值
		DepositEnable bool `json:"depositEnable"`
		// 是否是默认网络
		IsDefault bool `json:"isDefault"`
		// 备注正则
		MemoRegex string `json:"memoRegex"`
		// 最小确认数
		MinConfirm int `json:"minConfirm"`
		// 资产名
		Name string `json:"name"`
		// 网络
		Network string `json:"network"`
		// 特殊提示
		SpecialTips string `json:"specialTips"`
		// 解锁需要的确认数
		UnLockConfirm int `json:"unLockConfirm"`
		// 提现描述
		WithdrawDesc string `json:"withdrawDesc"`
		// 是否可以提现
		WithdrawEnable bool `json:"withdrawEnable"`
		// 提现手续费
		WithdrawFee string `json:"withdrawFee"`
		// 提现最小数量
		WithdrawMin string `json:"withdrawMin"`
		// 提现最大数量
		WithdrawMax string `json:"withdrawMax"`
		// 内部转账最小提现数
		WithdrawInternalMin string `json:"withdrawInternalMin"`
		// 是否需要memo
		SameAddress bool `json:"sameAddress"`
		// 预计到达时间
		EstimatedArrivalTime int `json:"estimatedArrivalTime"`
		// 是否繁忙
		Busy bool `json:"busy"`
		// 合约地址URL
		ContractAddressUrl string `json:"contractAddressUrl"`
		// 合约地址
		ContractAddress string `json:"contractAddress"`
		// 面额
		Denomination int `json:"denomination"`
	} `json:"networkList"`
}
