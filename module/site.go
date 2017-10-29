package module

import (
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type (
	Site struct{}

	siteState struct {
		AvatarURL string `json:"avatar_url"`
	}
)

// 站点状态
func (*Site) State(conf *Config) Value {
	state := &siteState{
		AvatarURL: com.AvatarURL("{name}", conf.Store.Url),
	}

	return &Succ{Data: state}
}
