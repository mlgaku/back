package conf

import (
	"github.com/mlgaku/back/module"
)

var Route = map[string]interface{}{
	"sub": &module.Sub{},

	"site": &module.Site{},
	"user": &module.User{},
	"bill": &module.Bill{},
	"node": &module.Node{},

	"topic": &module.Topic{},
	"reply": &module.Reply{},

	"notice": &module.Notice{},
}
