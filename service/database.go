package service

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Database struct {
	Session *mgo.Session
}

// 建立连接
func (d *Database) connect() {
	ses, err := mgo.Dial(fmt.Sprintf("%s:%d", APP.Conf.Db.Host, APP.Conf.Db.Port))
	if err != nil {
		panic(err)
	}

	ses.SetMode(mgo.Monotonic, true)
	d.Session = ses
}

// 获得 Collection 实例
func (d *Database) C(name string) *mgo.Collection {
	return d.Session.DB(APP.Conf.Db.Database).C(name)
}

// 通过 ID 查询
func (d *Database) FindId(tab, id string) *mgo.Query {
	return d.C(tab).FindId(bson.ObjectIdHex(id))
}

// 获得 Database 实例
func NewDatabase() *Database {
	db := &Database{}
	db.connect()
	return db
}
