package service

import (
	"gopkg.in/mgo.v2"
)

type Di struct {
	Module *Module
}

// 数据集合
func (*Di) C(n string) *mgo.Collection {
	return APP.Db.C(n)
}

// 发布订阅
func (*Di) Ps() *Pubsub {
	return APP.Ps
}

// 数据库
func (*Di) Db() *Database {
	return APP.Db
}

// 配置
func (*Di) Conf() *Config {
	return APP.Conf
}

// 模块
func (d *Di) Mod() *Module {
	return d.Module
}

// 会话
func (d *Di) Ses() *Session {
	return NewSession(d.Module.cli.Connection)
}

// 请求
func (d *Di) Req() *Request {
	return NewRequest([]byte(d.Module.Prot.Body), d.Module.cli)
}

// 响应
func (d *Di) Res() *Response {
	return NewResponse(d.Module.cli)
}
