package db

import (
	"errors"
	com "github.com/mlgaku/back/common"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	RegIP    string        `json:"reg_ip,omitempty" bson:"reg_ip"`
	RegDate  time.Time     `json:"reg_date,omitempty" bson:"reg_date"`
	Avatar   bool          `json:"avatar,omitempty" bson:",omitempty"`
	Balance  int64         `json:"balance"`
	Identity uint64        `json:"identity,omitempty" bson:",omitempty"`

	Name     string `fill:"i" json:"name" validate:"required,min=4,max=15,alphanum"`
	Email    string `fill:"i" json:"email" validate:"required,min=8,max=30,email"`
	Password string `fill:"i" json:"password,omitempty" validate:"required,min=8,max=20,alphanum"`

	Intro   string `fill:"u" json:"intro,omitempty" bson:",omitempty" validate:"omitempty,min=5,max=100"`
	Tagline string `fill:"u" json:"tagline,omitempty" bson:",omitempty" validate:"omitempty,min=3,max=30"`
	Website string `fill:"u" json:"website,omitempty" bson:",omitempty" validate:"omitempty,min=3,max=30,url"`

	service.Di
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
func (u *User) Add(user *User, salt string) error {
	if err := com.NewVali().Struct(user); err != "" {
		return errors.New(err)
	}

	user.RegDate = time.Now()
	user.Password = com.Sha1(user.Password, salt)
	return u.Db().C("user").Insert(user)
}

// 递增
func (u *User) Inc(id bson.ObjectId, field string, num int64) error {
	if id == "" {
		return errors.New("未指定用户ID")
	}
	return u.Db().C("user").UpdateId(id, M{"$inc": M{field: num}})
}

// 查找
func (u *User) Find(id bson.ObjectId, user interface{}, field M) error {
	return u.Db().C("user").FindId(id).Select(field).One(user)
}

// 保存
func (u *User) Save(id bson.ObjectId, user *User) error {
	if id == "" {
		return errors.New("用户ID不能为空")
	}

	set, err := com.Extract(user, "u")
	if err != nil {
		return err
	}

	return u.Db().C("user").UpdateId(id, M{"$set": set})
}

// 更新
func (u *User) Update(id bson.ObjectId, user M) error {
	return u.Db().C("user").UpdateId(id, M{"$set": user})
}

// 通过用户名查找
func (u *User) FindByName(name string, field M) (*User, error) {
	if name == "" {
		return nil, errors.New("用户名不能为空")
	}

	user := &User{}
	if err := u.Db().C("user").Find(M{"name": name}).Select(field).One(user); err != nil {
		return nil, errors.New(err.Error())
	}

	return user, nil
}

// 通过用户名查找多个
func (u *User) FindByNameMany(name []string) (map[string]User, error) {
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
	u.Db().C("user").Find(M{"name": M{"$in": in}}).All(&result)

	for _, v := range result {
		user[v.Name] = v
	}
	return user, nil
}

// 用户名是否存在
func (u *User) NameExists(name string) (bool, error) {
	if name == "" {
		return false, errors.New("用户名不能为空")
	}

	if c, _ := u.Db().C("user").Find(M{"name": name}).Count(); c > 0 {
		return true, nil
	}

	return false, nil
}

// 邮箱地址是否存在
func (u *User) EmailExists(email string) (bool, error) {
	if email == "" {
		return false, errors.New("邮箱地址不能为空")
	}

	if c, _ := u.Db().C("user").Find(M{"email": email}).Count(); c > 0 {
		return true, nil
	}

	return false, nil
}

// 通过ID修改头像状态
func (u *User) ChangeAvatarById(id bson.ObjectId, avatar bool) error {
	if id == "" {
		return errors.New("用户ID不能为空")
	}

	return u.Db().C("user").UpdateId(id, M{"$set": M{"avatar": avatar}})
}
