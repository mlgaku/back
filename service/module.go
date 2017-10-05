package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mlgaku/back/common"
	"github.com/mlgaku/back/conf"
	"github.com/mlgaku/back/types"
	"path"
	"reflect"
	"strings"
)

type (
	// 路由信息
	route struct {
		Mod  string `json:"mod"`  // 模块
		Act  string `json:"act"`  // 行为
		Body string `json:"body"` // 正文
	}
	// module 服务
	module struct {
		route    *route
		request  *request
		response *response
	}
)

// 加载模块
func (m *module) load(msg []byte) error {
	m.route = &route{}
	if json.Unmarshal(msg, m.route) != nil {
		return errors.New("json parsing failed")
	}

	switch {
	case m.route.Mod == "":
		return errors.New("mod get failed")
	case m.route.Act == "":
		return errors.New("act get failed")
	case !json.Valid([]byte(m.route.Body)):
		return errors.New("invalid body content")
	}

	m.request.body = m.route.Body
	return m.invoke()
}

// 打包数据
func (m *module) pack(data types.Value) []byte {
	m.route.Body = common.StringValue(&data)
	b, _ := json.Marshal(m.route)
	return b
}

// 调用方法
func (m *module) invoke() error {
	r, ok := conf.Route[m.route.Mod]
	if !ok {
		return fmt.Errorf("%s module does not exist", m.route.Mod)
	}

	mth := reflect.ValueOf(r).MethodByName(strings.Title(m.route.Act))
	if !mth.IsValid() {
		return fmt.Errorf("%s method does not exist", m.route.Act)
	}

	res := mth.Call(m.inject(&mth))
	if len(res) > 0 {
		m.response.write(m.pack(res[0].Interface()))
	}

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
		case "Database":
			args = append(args, reflect.ValueOf(_db.pseudo()))
		case "Config":
			args = append(args, reflect.ValueOf(_conf.pseudo()))

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
