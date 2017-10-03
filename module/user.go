package module

import (
	"encoding/json"
	"github.com/mlgaku/back/common"
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

	if err := common.NewValidator().Struct(user); err != "" {
		return &Fail{Msg: err}
	}

	if err := db.C("user").Insert(user); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{}
}
