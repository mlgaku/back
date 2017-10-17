package module

import (
	com "github.com/mlgaku/back/common"
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type User struct {
	Db db.User
}

// 注册
func (u *User) Reg(bd *Database, req *Request, conf *Config) Value {
	user, _ := db.NewUser(req.Body)

	user.RegIP, _ = com.IPAddr(req.RemoteAddr())
	if err := u.Db.Add(bd, conf, user); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{}
}

// 登录
func (u *User) Login(bd *Database, req *Request, ses *Session, conf *Config) Value {
	user, _ := db.NewUser(req.Body)
	if user.Password == "" {
		return &Fail{Msg: "密码不能为空"}
	}

	if _, ok := u.Check(bd, req).(*Succ); ok {
		return &Fail{Msg: "用户名不存在"}
	}

	result, err := u.Db.FindByName(bd, user.Name)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	if result.Password != com.Sha1(user.Password, conf.Secret.Salt) {
		return &Fail{Msg: "用户名与密码不匹配"}
	}

	// 保存状态
	ses.Set("user_id", result.Id)
	ses.Set("user_name", result.Name)
	ses.Set("user_email", result.Email)

	return &Succ{Data: &db.User{
		Id:    result.Id,
		Name:  result.Name,
		Email: result.Email,
	}}
}

// 检查用户名是否已被注册
func (u *User) Check(bd *Database, req *Request) Value {
	user, _ := db.NewUser(req.Body)

	b, err := u.Db.NameExists(bd, user.Name)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	if b {
		return &Fail{Msg: "用户名已存在"}
	}
	return &Succ{}
}
