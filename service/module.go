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
	app  *App
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
	return m.invoke()
}

// 调用方法
func (m *Module) invoke() (types.Value, error) {
	r, ok := m.app.Route[m.Prot.Mod]
	if !ok {
		return nil, fmt.Errorf("%s module does not exist", m.Prot.Mod)
	}

	mth := reflect.ValueOf(r).MethodByName(strings.Title(m.Prot.Act))
	if !mth.IsValid() {
		return nil, fmt.Errorf("%s method does not exist", m.Prot.Act)
	}

	res := mth.Call(m.inject(&mth))
	if len(res) > 0 {
		return res[0].Interface(), nil
		//m.response.Write(m.pack())
	}

	return nil, nil
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
		case "Request":
			args = append(args, reflect.ValueOf(NewRequest([]byte(m.Prot.Body), m.cli)))
		case "Response":
			args = append(args, reflect.ValueOf(NewResponse(m.cli)))
		case "Pubsub":
			args = append(args, reflect.ValueOf(m.app.Ps))
		case "Database":
			args = append(args, reflect.ValueOf(m.app.Db))
		case "Config":
			args = append(args, reflect.ValueOf(m.app.Conf))
		}
	}

	return args
}

// 获得 Module 实例
func NewModule(app *App, cli *Client) *Module {
	return &Module{
		app: app,
		cli: cli,
	}
}
