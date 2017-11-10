package module

import (
	"encoding/json"
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Topic struct {
	Db db.Topic
}

// 发表新主题
func (t *Topic) New(bd *Database, ses *Session, req *Request, conf *Config) Value {
	topic, _ := db.NewTopic(req.Body, "i")
	topic.Author = ses.Get("user").(*db.User).Id

	id, err := t.Db.Add(bd, topic)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	// 更新余额
	if conf.Reward.NewTopic != 0 {
		new(db.User).Inc(bd, topic.Author, "balance", conf.Reward.NewTopic)
		new(db.Bill).Add(bd, &db.Bill{
			Msg:    topic.Title,
			Type:   1,
			Date:   time.Now(),
			Number: conf.Reward.NewTopic,
			Master: topic.Author,
		})
	}

	return &Succ{Data: id}
}

// 编辑主题
func (t *Topic) Edit(ps *Pubsub, bd *Database, ses *Session, req *Request) Value {
	topic, _ := db.NewTopic(req.Body, "u")

	old := &db.Topic{}
	if err := t.Db.Find(bd, topic.Id, old); err != nil {
		return &Fail{Msg: err.Error()}
	}

	if err := t.Db.Save(bd, topic.Id, topic); err != nil {
		return &Fail{Msg: err.Error()}
	}

	// 添加通知
	typ := 0
	if topic.Title == old.Title && topic.Content == old.Content {
		if topic.Node != old.Node { // 移动
			typ = 4
		}
	} else {
		if topic.Node == old.Node { // 修改
			typ = 3
		} else { // 修改&移动
			typ = 5
		}
	}

	user := ses.Get("user").(*db.User)
	if typ != 0 && old.Author != user.Id {
		new(db.Notice).Add(bd, &db.Notice{
			Type:       uint64(typ),
			Date:       time.Now(),
			Master:     old.Author,
			User:       user.Name,
			TopicTitle: old.Title,
		})
	}

	ps.Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
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
	topic, _ := db.NewTopic(req.Body, "b")

	if err := t.Db.Find(bd, topic.Id, topic); err != nil {
		return &Fail{Msg: err.Error()}
	}

	t.Db.Inc(bd, topic.Id, "views")
	return &Succ{Data: topic}
}
