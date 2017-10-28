package conf

import (
	. "github.com/mlgaku/back/middleware"
)

var Middleware = map[string][]interface{}{
	// 发表主题
	"topic.new": {
		ShouldLogin,
	},

	// 回复主题
	"reply.new": {
		ShouldLogin,
	},

	// 获取通知
	"notice.list": {
		ShouldLogin,
	},

	// 添加节点
	"node.add": {
		ShouldLogin,
		ShouldFounder,
	},
	// 移除节点
	"node.remove": {
		ShouldLogin,
		ShouldFounder,
	},
}
