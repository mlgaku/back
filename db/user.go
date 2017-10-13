package db

import (
	"errors"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name" validate:"required,min=4,max=15,alphanum"`
	Email    string        `json:"email" validate:"required,min=8,max=30,email"`
	Password string        `json:"password,omitempty" validate:"required,min=8,max=20,alphanum"`

	RegIP   string    `json:"reg_ip,omitempty" bson:"reg_ip"`
	RegTime time.Time `json:"reg_time,omitempty" bson:"reg_time"`
}

// 添加
func (*User) Add(db *Database, conf *Config, user *User) error {
	if err := com.NewVali().Struct(user); err != "" {
		return errors.New(err)
	}

	user.RegTime = time.Now()
	user.Password = com.Sha1(user.Password, conf.Secret.Salt)
	if err := db.C("user").Insert(user); err != nil {
		return err
	}

	return nil
}

// 通过用户名查询
func (*User) FindByName(db *Database, name string) (*User, error) {
	if name == "" {
		return nil, errors.New("用户名不能为空")
	}

	user := &User{}
	if err := db.C("user").Find(bson.M{"name": name}).One(user); err != nil {
		return nil, errors.New(err.Error())
	}

	return user, nil
}

// 用户名是否存在
func (*User) NameExists(db *Database, name string) (bool, error) {
	if name == "" {
		return false, errors.New("用户名不能为空")
	}

	if c, _ := db.C("user").Find(bson.M{"name": name}).Count(); c > 0 {
		return true, nil
	}

	return false, nil
}
