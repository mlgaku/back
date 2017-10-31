package common

import (
	"encoding/json"
	. "github.com/mlgaku/back/types"
	"qiniupkg.com/x/errors.v7"
	"reflect"
	"strconv"
	"strings"
)

// 提取结构中的值
func Extract(v interface{}, str ...string) (map[string]interface{}, error) {
	val, result := reflect.ValueOf(v).Elem(), map[string]interface{}{}

	// 插入或更新模式
	if len(str) == 1 && len(str[0]) == 1 &&
		strings.Contains(ALLOW_INSERT+ALLOW_UPDATE, str[0]) {

		ele := reflect.TypeOf(v).Elem()
		for i, e := 0, ele.NumField(); i < e; i++ {
			f := ele.Field(i)

			// ID 不被更新
			if f.Name == "Id" && str[0] == ALLOW_UPDATE {
				continue
			}

			if strings.Contains(getFillType(f), str[0]) {
				w := val.Field(i).Interface()

				// 验证字段
				if v, ok := f.Tag.Lookup("validate"); ok {
					if err := NewVali().Var(w, v); err != "" {
						return nil, errors.New(err)
					}
				}

				result[getFieldName(f)] = w
			}
		}

		return result, nil
	}

	for i, e := 0, len(str); i < e; i++ {
		if f := val.FieldByName(strings.Title(str[i])); f.IsValid() {
			result[str[i]] = f.Interface()
		}
	}
	return result, nil
}

// 字符串值
func StringValue(val *Value) string {
	switch n := (*val).(type) {

	case int: // 数字
		return strconv.Itoa(n)

	case bool: // 布尔值
		if n {
			return "true"
		}
		return "false"

	case string: // 字符串
		return n

	case *Succ: // 成功
		n.Status = true

	case *Fail: // 失败
		n.Status = false
	}

	// 其它值统一转为json
	if b, e := json.Marshal(*val); e == nil {
		return string(b)
	}

	return "Conversion failed"
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

// 获取结构标签
func StructTag(v interface{}, field string, tag string) string {
	f, _ := reflect.TypeOf(v).Elem().FieldByName(strings.Title(field))
	return f.Tag.Get(tag)
}

// 获取填充类型
func getFillType(s reflect.StructField) string {
	return s.Tag.Get("fill")
}

// 获取字段名
func getFieldName(s reflect.StructField) string {
	key, ok := s.Tag.Lookup("bson")
	if ok {
		key = key[0:strings.Index(key, ",")]
	}

	if key == "" {
		key = strings.ToLower(s.Name)
	}

	return key
}
