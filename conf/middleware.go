package conf

import (
	mw "github.com/mlgaku/back/middleware"
)

var Middleware = map[string][]interface{}{
	"topic.new": {
		mw.IsLogin(),
	},
}
