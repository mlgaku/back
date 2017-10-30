package common

import (
	"encoding/json"
	"reflect"
	"strings"
)

const (
	ALLOW_BOTH   = "b"
	ALLOW_INSERT = "i"
	ALLOW_UPDATE = "u"
)

// 解析JSON
func ParseJSON(data []byte, typ string, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	if typ == ALLOW_BOTH {
		return nil
	}

	ele, val := reflect.TypeOf(v).Elem(), reflect.ValueOf(v).Elem()

	for i, e := 0, ele.NumField(); i < e; i++ {
		fill := ele.Field(i).Tag.Get("fill")
		if strings.Contains(fill, ALLOW_BOTH) {
			continue
		}
		if !strings.Contains(fill, typ) {
			val.Field(i).Set(reflect.Zero(val.Field(i).Type()))
		}
	}
	return nil
}
