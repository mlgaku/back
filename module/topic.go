package module

import (
	"encoding/json"
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
)

type Topic struct {
	Db db.Topic
}

// 发表新主题
func (t *Topic) New(bd *Database, ses *Session, req *Request) Value {
	topic, _ := db.NewTopic(req.Body)
	topic.Author = ses.Get("user").(db.User).Id

	id, err := t.Db.Add(bd, topic)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: id}
}

// 主题列表
func (t *Topic) List(bd *Database, req *Request) Value {
	var s struct {
		Page int
		Node bson.ObjectId
	}
	if err := json.Unmarshal(req.Body, &s); err != nil {
		return &Fail{Msg: err.Error()}
	}

	topic, err := t.Db.Paginate(bd, s.Node, s.Page)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: topic}
}

// 主题信息
func (t *Topic) Info(bd *Database, req *Request) Value {
	topic, _ := db.NewTopic(req.Body)

	if err := t.Db.Find(bd, topic.Id, topic); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: topic}
}
