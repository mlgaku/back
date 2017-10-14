package module

import (
	"encoding/json"
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
)

type Replay struct {
	Db db.Replay
}

func (*Replay) parse(body []byte) (*db.Replay, error) {
	replay := &db.Replay{}
	return replay, json.Unmarshal(body, replay)
}

// 添加新回复
func (r *Replay) New(db *Database, ps *Pubsub, ses *Session, req *Request) Value {
	replay, _ := r.parse(req.Body)
	replay.Author = ses.Get("user_id").(bson.ObjectId)
	replay.AuthorName = ses.Get("user_name").(string)

	if err := r.Db.Add(db, replay); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "replay", Act: "list"})
	ps.Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
}

// 获取回复列表
func (r *Replay) List(db *Database, req *Request) Value {
	var s struct {
		Page  int
		Topic bson.ObjectId
	}
	if err := json.Unmarshal(req.Body, &s); err != nil {
		return &Fail{Msg: err.Error()}
	}

	replay, err := r.Db.Paginate(db, s.Topic, s.Page)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: replay}
}
