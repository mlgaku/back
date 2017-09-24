package common

// 值
type Value interface{}

// 字符串值
func StringValue(v *Value) string {
	switch n := (*v).(type) {

	// 字符串
	case string:
		return n

	}

	return ""
}

// 字节值
func BytesValue(v *Value) []byte {
	return []byte(StringValue(v))
}
