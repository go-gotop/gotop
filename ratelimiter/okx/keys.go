package okx

const (
	// ============ 下单限制规则Key ================================
	// OkxCreateOrder2sKey 下单每2s不超过60次（单币种、期权除外）
	OkxCreateOrder2sKey = "okx:createorder:2s"
	// OkxCancelOrder2sKey 撤单每2s不超过60次
	OkxCancelOrder2sKey = "okx:cancelorder:2s"
)
