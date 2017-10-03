package service

import (
	"github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2"
)

type database struct {
	session *mgo.Session
}

// 建立连接
func (d *database) connect() {
	ses, err := mgo.Dial("127.0.0.1")
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
	return &types.Database{
		Session: d.session,
	}
}

// 获得 database 实例
func newDatabase() *database {
	db := &database{}
	db.connect()
	return db
}
