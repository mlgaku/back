package middleware

import (
	"errors"
	. "github.com/mlgaku/back/service"
)

// 登录
func IsLogin() interface{} {
	return func(ses *Session) error {
		if ses.Has("user_id") {
			return nil
		}
		return errors.New("你还没有登录")
	}
}

// 创始人
func IsFounder() interface{} {
	return func() error {
		return nil
	}
}

// 版主
func IsModerator() interface{} {
	return func() error {
		return nil
	}
}
