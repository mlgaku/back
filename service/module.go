package service

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/sxyazi/maile/conf"
	"log"
	"reflect"
	"strings"
)

// module 服务
type module struct {
	response *response
}

// 初始化
func (m *module) init(r *response) *module {
	m.response = r
	return m
}

// 加载模块
func (m *module) load(msg []byte) error {
	log.Println(string(msg))

	json, err := simplejson.NewJson(msg)
	if err != nil {
		return errors.New("json parsing failed.")
	}

	mod, err := json.Get("mod").String()
	if err != nil {
		return errors.New("mod get failed.")
	}

	act, err := json.Get("act").String()
	if err != nil {
		return errors.New("act get failed.")
	}

	return m.invoke(mod, act)
}

// 调用方法
func (m *module) invoke(mod, act string) error {
	r, ok := conf.Route[mod]
	if !ok {
		return fmt.Errorf("%s module does not exist.", mod)
	}

	res := reflect.ValueOf(r).MethodByName(strings.Title(act)).Call(nil)

	m.response.write(res[0].Interface())
	return nil
}
