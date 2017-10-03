package common

import (
	"encoding/json"
)

// 值
type Value interface{}

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

	// 其它值统一转为json
	default:
		b, e := json.Marshal(*val)
		if e == nil {
			return string(b)
		}

	}

	return "Conversion failed"
}

// 字节值
func BytesValue(val *Value) []byte {
	return []byte(StringValue(val))
}
