package module

import (
	"github.com/mlgaku/back/db"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Bill struct {
	db  db.Bill
	com common

	service.Di
}

// 获取账单列表
func (b *Bill) List() Value {
	dat, err := b.db.FindByMaster(b.com.user().Id)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: dat}
}
