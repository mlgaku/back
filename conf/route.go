package conf

import (
	"github.com/mlgaku/back/module"
)

var Route = map[string]interface{}{
	"sub": &module.Sub{},

	"home": &module.Home{},
	"user": &module.User{},
	"node": &module.Node{},
}
