package service

import (
	"encoding/json"
	"github.com/mlgaku/back/types"
	"io/ioutil"
)

type config struct {
	conf *types.Config
}

// 读配置
func (c *config) read() {
	env, err := ioutil.ReadFile(".env")
	if err != nil {
		panic(err)
	}
	c.conf = &types.Config{}
	if err = json.Unmarshal(env, c.conf); err != nil {
		panic(err)
	}
}

// 创建替身
func (c *config) pseudo() *types.Config {
	return c.conf
}

// 获得 config 实例
func newConfig() *config {
	c := &config{}
	c.read()
	return c
}
