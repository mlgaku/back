package common

// 值
type Value interface{}

// 字符串值
func StringValue(val *Value) string {
	switch n := (*val).(type) {

	// 字符串
	case string:
		return n

	}

	return ""
}

// 字节值
func BytesValue(val *Value) []byte {
	return []byte(StringValue(val))
}
