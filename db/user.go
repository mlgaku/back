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
	RegIP    string        `json:"reg_ip,omitempty" bson:"reg_ip"`
	RegDate  time.Time     `json:"reg_date,omitempty" bson:"reg_date"`
	Avatar   bool          `json:"avatar,omitempty" bson:",omitempty"`
	Identity uint64        `json:"identity,omitempty" bson:",omitempty"`

	Name     string `fill:"i" json:"name" validate:"required,min=4,max=15,alphanum"`
	Email    string `fill:"i" json:"email" validate:"required,min=8,max=30,email"`
	Password string `fill:"i" json:"password,omitempty" validate:"required,min=8,max=20,alphanum"`

	Intro   string `fill:"u" json:"intro,omitempty" bson:",omitempty" validate:"omitempty,min=5,max=100"`
	Tagline string `fill:"u" json:"tagline,omitempty" bson:",omitempty" validate:"omitempty,min=3,max=30"`
	Website string `fill:"u" json:"website,omitempty" bson:",omitempty" validate:"omitempty,min=3,max=30,url"`
}

// 获得 User 实例
func NewUser(body []byte, typ string) (*User, error) {
	user := &User{}
	if err := com.ParseJSON(body, typ, user); err != nil {
		panic(err)
	}

	return user, nil
}

// 添加
func (*User) Add(db *Database, conf *Config, user *User) error {
	if err := com.NewVali().Struct(user); err != "" {
		return errors.New(err)
	}

	user.RegDate = time.Now()
	user.Password = com.Sha1(user.Password, conf.Secret.Salt)
	return db.C("user").Insert(user)
}

// 查找
func (*User) Find(db *Database, id bson.ObjectId, user interface{}, field bson.M) error {
	return db.C("user").FindId(id).Select(field).One(user)
}

// 保存
func (*User) Save(db *Database, id bson.ObjectId, user *User) error {
	if id == "" {
		return errors.New("用户ID不能为空")
	}

	set, err := com.Extract(user, "u")
	if err != nil {
		return err
	}

	return db.C("user").UpdateId(id, bson.M{"$set": set})
}

// 更新
func (*User) Update(db *Database, id bson.ObjectId, user bson.M) error {
	return db.C("user").UpdateId(id, bson.M{"$set": user})
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

// 通过ID修改头像状态
func (*User) ChangeAvatarById(db *Database, id bson.ObjectId, avatar bool) error {
	if id == "" {
		return errors.New("用户ID不能为空")
	}

	return db.C("user").UpdateId(id, bson.M{"$set": bson.M{"avatar": avatar}})
}
