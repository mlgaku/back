package conf

import (
	. "github.com/mlgaku/back/middleware"
)

var Middleware = map[string][]interface{}{
	// 用户
	"user.info":           {ShouldLogin},
	"user.avatar":         {ShouldLogin},
	"user.setAvatar":      {ShouldLogin},
	"user.removeAvatar":   {ShouldLogin},
	"user.changePassword": {ShouldLogin},

	// 主题
	"topic.new": {ShouldLogin},

	// 回复
	"reply.new": {ShouldLogin},

	// 通知
	"notice.list": {ShouldLogin},

	// 节点
	"node.add":    {ShouldLogin, ShouldFounder},
	"node.remove": {ShouldLogin, ShouldFounder},
}
