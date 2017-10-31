package module

import (
	com "github.com/mlgaku/back/common"
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type User struct {
	Db db.User
}

// 注册
func (u *User) Reg(bd *Database, req *Request, conf *Config) Value {
	user, _ := db.NewUser(req.Body, "i")

	user.RegIP, _ = com.IPAddr(req.RemoteAddr())
	if err := u.Db.Add(bd, conf, user); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{}
}

// 登录
func (u *User) Login(bd *Database, req *Request, ses *Session, conf *Config) Value {
	user, _ := db.NewUser(req.Body, "b")
	if user.Password == "" {
		return &Fail{Msg: "密码不能为空"}
	}

	if v, ok := u.Check(bd, req).(*Succ); !ok {
		return &Fail{Msg: "检查用户名失败"}
	} else if v.Data == true {
		return &Fail{Msg: "用户名不存在"}
	}

	result, err := u.Db.FindByName(bd, user.Name)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	if result.Password != com.Sha1(user.Password, conf.Secret.Salt) {
		return &Fail{Msg: "用户名与密码不匹配"}
	}

	ses.Set("user", result)
	return &Succ{}
}

// 用户信息
func (u *User) Info(bd *Database, ses *Session, conf *Config) Value {
	user := ses.Get("user").(*db.User)

	result := &db.User{}
	if err := u.Db.Find(bd, user.Id, result); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: result}
}

// 检查用户名是否已被注册
func (u *User) Check(bd *Database, req *Request) Value {
	user, _ := db.NewUser(req.Body, "b")

	b, err := u.Db.NameExists(bd, user.Name)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: !b}
}

// 检查邮箱地址是否已存在
func (u *User) CheckEmail(bd *Database, req *Request) Value {
	user, _ := db.NewUser(req.Body, "b")

	b, err := u.Db.EmailExists(bd, user.Email)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: !b}
}

// 上传头像
func (u *User) Avatar(ses *Session, conf *Config) Value {
	file := com.AvatarFile(ses.Get("user").(*db.User).Name)

	policy := storage.PutPolicy{
		Expires:    120,
		DetectMime: 1,
		FsizeLimit: 1048576,
		MimeLimit:  "image/*",
		Scope:      conf.Store.Bucket + ":" + file,
	}

	mac := qbox.NewMac(conf.Store.Ak, conf.Store.Sk)
	return &Succ{Data: map[string]string{
		"file":  file,
		"token": policy.UploadToken(mac),
	}}
}

// 设置头像
func (u *User) SetAvatar(ps *Pubsub, bd *Database, ses *Session, conf *Config) Value {
	user := ses.Get("user").(*db.User)

	if err := u.Db.ChangeAvatarById(bd, user.Id, true); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "user", Act: "info"})
	return &Succ{}
}

// 移除头像
func (u *User) RemoveAvatar(ps *Pubsub, bd *Database, ses *Session, conf *Config) Value {
	user := ses.Get("user").(*db.User)

	// 删除头像文件
	manager := storage.NewBucketManager(qbox.NewMac(conf.Store.Ak, conf.Store.Sk), nil)
	if err := manager.Delete(conf.Store.Bucket, com.AvatarFile(user.Name)); err != nil {
		return &Fail{Msg: err.Error()}
	}

	// 改变头像状态
	if err := u.Db.ChangeAvatarById(bd, user.Id, false); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "user", Act: "info"})
	return &Succ{}
}

// 编辑资料
func (u *User) EditProfile(bd *Database, ses *Session, req *Request) Value {
	user, _ := db.NewUser(req.Body, "u")

	if err := u.Db.Save(bd, ses.Get("user").(*db.User).Id, user); err != nil {
		return &Fail{Msg: err.Error()}
	}
	return &Succ{}
}

// 更改密码
func (u *User) ChangePassword(conf *Config) {

}
