package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sxyazi/maile/common"
	"github.com/sxyazi/maile/conf"
	"path"
	"reflect"
	"strings"
)

// module 服务
type module struct {
	request  *request
	response *response
}

// 加载模块
func (m *module) load(msg []byte) error {
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

	res := mth.Call(m.inject(&mth))
	m.response.write(res[0].Interface())
	return nil
}

// 依赖注入
func (m *module) inject(mth *reflect.Value) []reflect.Value {
	num := (*mth).Type().NumIn()
	if num < 1 {
		return nil
	}

	args := make([]reflect.Value, 0)
	for i := 0; i < num; i++ {
		t := (*mth).Type().In(i)
		if t.Kind() != reflect.Ptr {
			args = append(args, reflect.ValueOf(nil))
			continue
		}

		switch n := strings.TrimLeft(path.Ext(t.String()), "."); n {
		case "Request":
			args = append(args, reflect.ValueOf(m.request.pseudo()))
		case "Response":
			args = append(args, reflect.ValueOf(m.response.pseudo()))
		}
	}

	return args
}

// 获得 module 实例
func newModule(req *request, res *response) *module {
	return &module{
		request:  req,
		response: res,
	}
}
