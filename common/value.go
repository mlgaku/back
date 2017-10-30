package common

import (
	"encoding/json"
	. "github.com/mlgaku/back/types"
	"reflect"
	"strings"
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

// 过滤结构中的空值
func FilterStruct(v interface{}) map[string]interface{} {
	ele, val := reflect.TypeOf(v).Elem(), reflect.ValueOf(v).Elem()

	result := map[string]interface{}{}
	for i, e := 0, val.NumField(); i < e; i++ {
		f := val.Field(i)

		if f.Interface() != reflect.Zero(f.Type()).Interface() {
			result[getFieldName(ele.Field(i))] = f.Interface()
		}
	}

	return result
}

// 获取字段名
func getFieldName(s reflect.StructField) string {
	key, ok := s.Tag.Lookup("bson")
	if ok {
		key = key[0:strings.Index(key, ",")]
		if key == "" {
			key = strings.ToLower(s.Name)
		}
	}

	if key == "" {
		return s.Name
	}

	return key
}
