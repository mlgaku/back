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
	"user.editProfile":    {ShouldLogin},
	"user.changePassword": {ShouldLogin},

	// 账单
	"bill.list": {ShouldLogin},

	// 主题
	"topic.new":    {ShouldLogin},
	"topic.edit":   {ShouldLogin, ShouldModerator},
	"topic.subtle": {ShouldLogin},

	// 回复
	"reply.new":  {ShouldLogin},
	"reply.edit": {ShouldLogin, ShouldModerator},

	// 通知
	"notice.list": {ShouldLogin},

	// 节点
	"node.add":    {ShouldLogin, ShouldFounder},
	"node.edit":   {ShouldLogin, ShouldFounder},
	"node.remove": {ShouldLogin, ShouldFounder},
}
