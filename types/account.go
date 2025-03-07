package types

import (
	"fmt"
	"strings"
)

// AccountType 账户类型：1-AccountTypeClassic, 2-AccountTypeUnified
type AccountType int

// String 返回字符串表示
func (a AccountType) String() string {
	switch a {
	case AccountTypeClassic:
		return "CLASSIC"
	case AccountTypeUnified:
		return "UNIFIED"
	default:
		return "UNKNOWN"
	}
}

// IsValid 判断 AccountType 是否为已定义的类型
func (a AccountType) IsValid() bool {
	switch a {
	case AccountTypeClassic, AccountTypeUnified:
		return true
	default:
		return false
	}
}

// ParseAccountType 从字符串解析 AccountType (不区分大小写)
func ParseAccountType(s string) (AccountType, error) {
	s = strings.ToUpper(strings.TrimSpace(s))
	switch s {
	case "CLASSIC":
		return AccountTypeClassic, nil
	case "UNIFIED":
		return AccountTypeUnified, nil
	default:
		return AccountTypeUnknown, fmt.Errorf("unknown account type: %s", s)
	}
}

const (
	// AccountTypeUnknown 未知账户
	AccountTypeUnknown AccountType = iota
	// AccountTypeClassic 经典账户
	AccountTypeClassic
	// AccountTypeUnified 统一账户
	AccountTypeUnified
)

type Account struct {
}
