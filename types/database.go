package types

import "gopkg.in/mgo.v2"

type Database struct {
	Config struct {
		Host string // 主机
		Port int    // 端口
		Name string // 数据库名
	}
	Session *mgo.Session
}

func (d *Database) C(name string) *mgo.Collection {
	return d.Session.DB(d.Config.Name).C(name)
}
