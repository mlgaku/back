package module

import (
	"github.com/mlgaku/back/db"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Bill struct {
	db db.Bill

	service.Di
}

// 获取账单列表
func (b *Bill) List() Value {
	return &Succ{
		Data: b.db.FindByMaster(b.Ses().Get("user").(*db.User).Id),
	}
}
