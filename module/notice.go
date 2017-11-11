package module

import (
	"github.com/mlgaku/back/db"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Notice struct {
	db db.Notice

	service.Di
}

// 获取通知列表
func (n *Notice) List() Value {
	dat, err := n.db.FindByMaster(n.Ses().Get("user").(*db.User).Id)

	if err != nil {
		return &Fail{Msg: err.Error()}
	}
	return &Succ{Data: dat}
}

// 移除通知
func (n *Notice) Remove() Value {
	notice, _ := db.NewNotice(n.Req().Body, "b")

	if err := n.db.Find(notice.Id, notice); err != nil {
		return &Fail{Msg: err.Error()}
	}

	if notice.Master != n.Ses().Get("user").(*db.User).Id {
		return &Fail{Msg: "你不能移除别人的通知"}
	}

	if err := n.db.ChangeReadById(notice.Id, true); err != nil {
		return &Fail{Msg: err.Error()}
	}

	n.Ps().Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
}
