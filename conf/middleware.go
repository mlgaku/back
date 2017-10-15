package conf

import (
	. "github.com/mlgaku/back/middleware"
)

var Middleware = map[string][]interface{}{
	"topic.new": {
		ShouldLogin,
	},
	"reply.new": {
		ShouldLogin,
	},
}
