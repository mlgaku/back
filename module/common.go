package module

import (
	"github.com/mlgaku/back/db"
	"github.com/mlgaku/back/service"
)

type common struct {
	service.Di
}

// 获取用户 Session
func (c *common) user() *db.User {
	return c.Ses().Get("user").(*db.User)
}
