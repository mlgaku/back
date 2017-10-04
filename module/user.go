package module

import (
	"encoding/json"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/types"
)

type User struct {
	Name     string `json:"name" validate:"required,min=4,max=15,alphanum"`
	Email    string `json:"email" validate:"required,min=8,max=30,email"`
	Password string `json:"password" validate:"required,min=8,max=20,alphanum"`
}

// 注册
func (*User) Reg(db *Database, req *Request) Value {

	user := &User{}
	json.Unmarshal(req.Body, user)

	if err := com.NewVali().Struct(user); err != "" {
		return &Fail{Msg: err}
	}

	if err := db.C("user").Insert(user); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{}
}

// 登录
func (*User) Login(db *Database, req *Request) Value {

	user := &User{}
	json.Unmarshal(req.Body, user)

	err := com.NewVali().Each(com.Iter(user.Name, user.Password), []string{"required"})
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	if c, _ := db.C("user").Find(user).Count(); c != 1 {
		return &Fail{Msg: "登录失败"}
	}

	return &Succ{}
}
