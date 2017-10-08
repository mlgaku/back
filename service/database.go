package service

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

type Database struct {
	config struct {
		name, host string
		port       int
	}
	Session *mgo.Session
}

// 建立连接
func (d *Database) connect() {
	ses, err := mgo.Dial(fmt.Sprintf("%s:%d", d.config.host, d.config.port))
	if err != nil {
		panic(err)
	}

	ses.SetMode(mgo.Monotonic, true)
	d.Session = ses
}

// 获得 Collection 实例
func (d *Database) C(name string) *mgo.Collection {
	return d.Session.DB(d.config.name).C(name)
}

// 获得 Database 实例
func NewDatabase(host, name string, port int) *Database {
	db := &Database{}
	db.config.name, db.config.host, db.config.port = name, host, port

	db.connect()
	return db
}
