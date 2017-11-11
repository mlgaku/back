package module

import (
	"encoding/json"
	com "github.com/mlgaku/back/common"
	"github.com/mlgaku/back/db"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type (
	User struct {
		db db.User

		service.Di
	}

	// 用户主页
	userHome struct {
		User  *db.User    `json:"user"`
		Topic *[]db.Topic `json:"topic"`
		Reply *[]db.Reply `json:"reply"`
	}
)

// 注册
func (u *User) Reg() Value {
	user, _ := db.NewUser(u.Req().Body, "i")

	user.RegIP, _ = com.IPAddr(u.Req().RemoteAddr())
	if err := u.db.Add(user, u.Conf().Secret.Salt); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{}
}

// 登录
func (u *User) Login() Value {
	user, _ := db.NewUser(u.Req().Body, "b")
	if user.Password == "" {
		return &Fail{Msg: "密码不能为空"}
	}

	if v, ok := u.Check().(*Succ); !ok {
		return &Fail{Msg: "检查用户名失败"}
	} else if v.Data == true {
		return &Fail{Msg: "用户名不存在"}
	}

	result, err := u.db.FindByName(user.Name, nil)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	if result.Password != com.Sha1(user.Password, u.Conf().Secret.Salt) {
		return &Fail{Msg: "用户名与密码不匹配"}
	}

	u.Ses().Set("user", result)
	return &Succ{}
}

// 用户主页
func (u *User) Home() Value {
	user, _ := db.NewUser(u.Req().Body, "b")
	home := &userHome{}

	err := error(nil)
	if home.User, err = u.db.FindByName(user.Name, M{
		"reg_ip":   0,
		"password": 0,
	}); err != nil {
		return &Fail{Msg: err.Error()}
	}

	if home.Topic, err = new(db.Topic).FindByAuthor(home.User.Id, M{"content": 0}, 0); err != nil {
		return &Fail{Msg: err.Error()}
	}

	if home.Reply, err = new(db.Reply).FindByAuthor(home.User.Id, nil, 0); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: home}
}

// 用户信息
func (u *User) Info() Value {
	user := u.Ses().Get("user").(*db.User)

	result := &db.User{}
	if err := u.db.Find(user.Id, result, M{
		"reg_ip":   0,
		"password": 0,
	}); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: result}
}

// 检查用户名是否已被注册
func (u *User) Check() Value {
	user, _ := db.NewUser(u.Req().Body, "b")

	b, err := u.db.NameExists(user.Name)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: !b}
}

// 检查邮箱地址是否已存在
func (u *User) CheckEmail() Value {
	user, _ := db.NewUser(u.Req().Body, "b")

	b, err := u.db.EmailExists(user.Email)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: !b}
}

// 上传头像
func (u *User) Avatar() Value {
	conf, file := u.Conf(), com.AvatarFile(u.Ses().Get("user").(*db.User).Name)

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
func (u *User) SetAvatar() Value {
	user := u.Ses().Get("user").(*db.User)

	if err := u.db.ChangeAvatarById(user.Id, true); err != nil {
		return &Fail{Msg: err.Error()}
	}

	u.Ps().Publish(&Prot{Mod: "user", Act: "info"})
	return &Succ{}
}

// 移除头像
func (u *User) RemoveAvatar() Value {
	conf, user := u.Conf(), u.Ses().Get("user").(*db.User)

	// 删除头像文件
	manager := storage.NewBucketManager(qbox.NewMac(conf.Store.Ak, conf.Store.Sk), nil)
	if err := manager.Delete(conf.Store.Bucket, com.AvatarFile(user.Name)); err != nil {
		return &Fail{Msg: err.Error()}
	}

	// 改变头像状态
	if err := u.db.ChangeAvatarById(user.Id, false); err != nil {
		return &Fail{Msg: err.Error()}
	}

	u.Ps().Publish(&Prot{Mod: "user", Act: "info"})
	return &Succ{}
}

// 编辑资料
func (u *User) EditProfile() Value {
	user, _ := db.NewUser(u.Req().Body, "u")

	if err := u.db.Save(u.Ses().Get("user").(*db.User).Id, user); err != nil {
		return &Fail{Msg: err.Error()}
	}

	u.Ps().Publish(&Prot{Mod: "user", Act: "info"})
	return &Succ{}
}

// 更改密码
func (u *User) ChangePassword() Value {
	var j struct {
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}

	json.Unmarshal(u.Req().Body, &j)
	if j.Password == j.NewPassword {
		return &Succ{}
	}

	if err := com.NewVali().Var(
		j.NewPassword,
		com.StructTag(&u.db, "password", "validate"),
	); err != "" {
		return &Fail{Msg: err}
	}

	id, conf := u.Ses().Get("user").(*db.User).Id, u.Conf()

	user := &db.User{}
	u.db.Find(id, user, M{"password": 1})
	if com.Sha1(j.Password, conf.Secret.Salt) != user.Password {
		return &Fail{Msg: "原密码输入不正确"}
	}

	u.db.Update(id, M{"password": com.Sha1(j.NewPassword, conf.Secret.Salt)})
	return &Succ{}
}
