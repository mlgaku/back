package middleware

import (
	"errors"
	. "github.com/mlgaku/back/service"
)

type Auth struct{}

func (*Auth) IsLogin() interface{} {
	return func(ses *Session) error {
		if ses.Has("user_id") {
			return nil
		}
		return errors.New("你还没有登录")
	}
}
