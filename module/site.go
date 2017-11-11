package module

import (
	com "github.com/mlgaku/back/common"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Site struct {
	service.Di
}

// 站点状态
func (s *Site) State() Value {
	return &Succ{Data: map[string]string{
		"avatar_url": com.AvatarURL("{name}", s.Conf().Store.Url),
	}}
}
