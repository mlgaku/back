package db

import (
	"encoding/json"
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
	Identity uint64        `json:"identity,omitempty" bson:",omitempty"`

	RegIP   string    `json:"reg_ip,omitempty" bson:"reg_ip"`
	RegDate time.Time `json:"reg_date,omitempty" bson:"reg_date"`
}

// 获得 User 实例
func NewUser(body []byte) (*User, error) {
	user := &User{}
	if err := json.Unmarshal(body, user); err != nil {
		return nil, err
	}

	user.Identity, user.RegIP, user.RegDate = 0, "", time.Now()
	return user, nil
}

// 添加
func (*User) Add(db *Database, conf *Config, user *User) error {
	if err := com.NewVali().Struct(user); err != "" {
		return errors.New(err)
	}

	user.Password = com.Sha1(user.Password, conf.Secret.Salt)
	return db.C("user").Insert(user)
}

// 查找
func (*User) Find(db *Database, id bson.ObjectId, user interface{}) error {
	return db.C("user").FindId(id).One(user)
}

// 通过用户名查找
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

// 通过用户名查找多个
func (*User) FindByNameMany(db *Database, name []string) (map[string]User, error) {
	if len(name) < 1 {
		return nil, nil
	}

	in, user := []string{}, map[string]User{}
	for _, v := range name {
		if _, ok := user[v]; !ok {
			user[v] = User{}
			in = append(in, v)
		}
	}

	result := []User{}
	db.C("user").Find(bson.M{"name": bson.M{"$in": in}}).All(&result)

	for _, v := range result {
		user[v.Name] = v
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

// 邮箱地址是否存在
func (*User) EmailExists(db *Database, email string) (bool, error) {
	if email == "" {
		return false, errors.New("邮箱地址不能为空")
	}

	if c, _ := db.C("user").Find(bson.M{"email": email}).Count(); c > 0 {
		return true, nil
	}

	return false, nil
}
