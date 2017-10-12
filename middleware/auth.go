package middleware

import (
	"errors"
	. "github.com/mlgaku/back/service"
)

// 登录
func ShouldLogin(ses *Session) error {
	if ses.Has("user_id") {
		return nil
	}
	return errors.New("你还没有登录")
}

// 创始人
func ShouldFounder() error {
	return nil
}

// 版主
func ShouldModerator() error {
	return nil
}
