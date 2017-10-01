package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sxyazi/maile/common"
	"github.com/sxyazi/maile/conf"
	"log"
	"reflect"
	"strings"
)

// module 服务
type module struct {
	response *response
}

// 加载模块
func (m *module) load(msg []byte) error {
	log.Println(string(msg))

	route := &common.Route{}
	if json.Unmarshal(msg, route) != nil {
		return errors.New("json parsing failed.")
	}

	switch {
	case route.Module == "":
		return errors.New("mod get failed.")
	case route.Action == "":
		return errors.New("act get failed.")
	}

	return m.invoke(route.Module, route.Action)
}

// 调用方法
func (m *module) invoke(mod, act string) error {
	r, ok := conf.Route[mod]
	if !ok {
		return fmt.Errorf("%s module does not exist.", mod)
	}

	mth := reflect.ValueOf(r).MethodByName(strings.Title(act))
	if !mth.IsValid() {
		return fmt.Errorf("%s method does not exist.", act)
	}

	res := mth.Call(nil)
	m.response.write(res[0].Interface())
	return nil
}

// 获得 module 实例
func newModule(res *response) *module {
	return &module{
		response: res,
	}
}
