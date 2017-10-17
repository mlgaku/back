package middleware

import (
	"errors"
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
)

// 登录
func ShouldLogin(ses *Session) error {
	if ses.Has("user") {
		return nil
	}
	return errors.New("你还没有登录")
}

// 创始人
func ShouldFounder(ses *Session) error {
	if ses.Get("user").(db.User).Identity == 2 {
		return nil
	}
	return errors.New("你不能执行此操作")
}

// 版主
func ShouldModerator(ses *Session) error {
	if ses.Get("user").(db.User).Identity >= 1 {
		return nil
	}
	return errors.New("你不能执行此操作")
}
