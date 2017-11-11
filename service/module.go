package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mlgaku/back/types"
	"log"
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
	if json.Unmarshal(msg, m.Prot) != nil {
		return nil, errors.New("json parsing failed")
	}

	return m.LoadProt(nil)
}

// 加载模块(Prot方式)
func (m *Module) LoadProt(prot *types.Prot) (types.Value, error) {
	if prot != nil {
		m.Prot = prot
	}

	switch {
	case m.Prot.Mod == "":
		return nil, errors.New("mod get failed")
	case m.Prot.Act == "":
		return nil, errors.New("act get failed")
	case !json.Valid([]byte(m.Prot.Body)):
		return nil, errors.New("invalid body content")
	}

	// 方法名统一小写
	m.Prot.Act = strings.ToLower(m.Prot.Act[:1]) + m.Prot.Act[1:]
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

	// 中间件
	if err := m.middle(); err != nil {
		return nil, err
	}

	// 填充当前 module 状态
	if di := reflect.ValueOf(r).Elem().FieldByName("Di"); di.IsValid() {
		if f := di.FieldByName("Module"); f.CanSet() {
			f.Set(reflect.ValueOf(m))
		}
	}

	if !APP.Conf.App.Debug {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}
		}()
	}

	// 调用方法
	if res := mth.Call(nil); len(res) > 0 {
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

// 依赖注入(中间件)
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
		cli:  cli,
		Prot: new(types.Prot),
	}
}
