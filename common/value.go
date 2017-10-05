package common

import (
	"encoding/json"
	. "github.com/mlgaku/back/types"
)

// 字符串值
func StringValue(val *Value) string {
	switch n := (*val).(type) {

	// 数字
	case int:
		return string(n)

	// 布尔值
	case bool:
		if n {
			return "true"
		}
		return "false"

	// 字符串
	case string:
		return n

	// 成功
	case *Succ:
		n.Status = true

	// 失败
	case *Fail:
		n.Status = false

	}

	// 其它值统一转为json
	if b, e := json.Marshal(*val); e == nil {
		return string(b)
	}

	return "Conversion failed"
}

// 字节值
func BytesValue(val *Value) []byte {
	return []byte(StringValue(val))
}
