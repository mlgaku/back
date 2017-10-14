package module

import (
	"encoding/json"
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Notice struct {
	Db db.Notice
}

func (*Notice) parse(body []byte) (*db.Notice, error) {
	replay := &db.Notice{}
	return replay, json.Unmarshal(body, replay)
}

// 获取通知列表
func (n *Notice) List(db *Database, req *Request) Value {
	notice, _ := n.parse(req.Body)

	dat, err := n.Db.FindByMaster(db, notice.Master)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: dat}
}

// 移除通知
func (n *Notice) Remove(ps *Pubsub, db *Database, req *Request) Value {
	notice, _ := n.parse(req.Body)

	if err := n.Db.ChangeReadById(db, notice.Id, true); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
}
