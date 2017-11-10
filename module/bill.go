package module

import (
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Bill struct {
	Db db.Bill
}

// 获取账单列表
func (b *Bill) List(bd *Database, ses *Session) Value {
	dat, err := b.Db.FindByMaster(bd, ses.Get("user").(*db.User).Id)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: dat}
}
