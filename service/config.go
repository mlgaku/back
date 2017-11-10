package service

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Db struct {
		Host     string
		Port     int
		Database string
	}
	App struct {
		Debug bool
	}
	Store struct {
		Ak     string
		Sk     string
		Url    string
		Bucket string
	}
	Secret struct {
		Salt string
	}
	Reward struct {
		Sign     int64
		NewTopic int64 `json:"new_topic"`
		NewReply int64 `json:"new_reply"`
	}
}

// 读配置
func (c *Config) read() {
	env, err := ioutil.ReadFile(".env")
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(env, c); err != nil {
		panic(err)
	}
}

// 获得 Config 实例
func NewConfig() *Config {
	c := &Config{}
	c.read()
	return c
}
