package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mlgaku/back/types"
	"path"
	"reflect"
	"strings"
)

type Module struct {
	cli  *Client
	Prot *types.Prot
}

// 加载模块
func (m *Module) Load(msg []byte) (types.Value, error) {
	p := &types.Prot{}
	if json.Unmarshal(msg, p) != nil {
		return nil, errors.New("json parsing failed")
	}

	return m.LoadProt(p)
}

// 加载模块(Prot方式)
func (m *Module) LoadProt(prot *types.Prot) (types.Value, error) {
	switch {
	case prot.Mod == "":
		return nil, errors.New("mod get failed")
	case prot.Act == "":
		return nil, errors.New("act get failed")
	case !json.Valid([]byte(prot.Body)):
		return nil, errors.New("invalid body content")
	}

	m.Prot = prot
	m.Prot.Act = strings.ToLower(m.Prot.Act)
	return m.invoke()
}

// 调用方法
func (m *Module) invoke() (types.Value, error) {
	r, ok := APP.Route[m.Prot.Mod]
	if !ok {
		return nil, fmt.Errorf("%s module does not exist", m.Prot.Mod)
	}

	mth := reflect.ValueOf(r).MethodByName(strings.Title(m.Prot.Act))
	if !mth.IsValid() {
		return nil, fmt.Errorf("%s method does not exist", m.Prot.Act)
	}

	if err := m.middle(); err != nil {
		return nil, err
	}

	res := mth.Call(m.inject(&mth))
	if len(res) > 0 {
		return res[0].Interface(), nil
	}

	return nil, nil
}

// 中间件
func (m *Module) middle() error {
	w, ok := APP.Middleware[m.Prot.Mod+"."+m.Prot.Act]
	if !ok {
		return nil
	}

	for _, v := range w {
		fn := reflect.ValueOf(v)
		res := fn.Call(m.inject(&fn))

		if len(res) < 1 {
			continue
		}

		if err, ok := res[0].Interface().(error); ok && err != nil {
			return err
		}
	}

	return nil
}

// 依赖注入
func (m *Module) inject(mth *reflect.Value) []reflect.Value {
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
		case "Module":
			args = append(args, reflect.ValueOf(m))
		case "Session":
			args = append(args, reflect.ValueOf(NewSession(m.cli.Connection)))
		case "Request":
			args = append(args, reflect.ValueOf(NewRequest([]byte(m.Prot.Body), m.cli)))
		case "Response":
			args = append(args, reflect.ValueOf(NewResponse(m.cli)))
		case "Pubsub":
			args = append(args, reflect.ValueOf(APP.Ps))
		case "Database":
			args = append(args, reflect.ValueOf(APP.Db))
		case "Config":
			args = append(args, reflect.ValueOf(APP.Conf))
		}
	}

	return args
}

// 获得 Module 实例
func NewModule(cli *Client) *Module {
	return &Module{
		cli: cli,
	}
}
