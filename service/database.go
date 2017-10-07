package service

import (
	"fmt"
	"github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2"
)

type database struct {
	config struct {
		name, host string
		port       int
	}
	session *mgo.Session
}

// 建立连接
func (d *database) connect() {
	ses, err := mgo.Dial(fmt.Sprintf("%s:%d", d.config.host, d.config.port))
	if err != nil {
		panic(err)
	}

	ses.SetMode(mgo.Monotonic, true)
	d.session = ses
}

// 断开连接
func (d *database) disconnect() {
	d.session.Close()
}

// 创建替身
func (d *database) pseudo() *types.Database {
	t := &types.Database{Session: d.session}
	t.Config.Host, t.Config.Name, t.Config.Port = d.config.host, d.config.name, d.config.port

	return t
}

// 获得 database 实例
func newDatabase(host, name string, port int) *database {
	db := &database{}
	db.config.name, db.config.host, db.config.port = name, host, port

	db.connect()
	return db
}
