package module

import (
	"encoding/json"
	"github.com/mlgaku/back/db"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Topic struct {
	Db db.Topic
	service.Di
}

// 发表新主题
func (t *Topic) New() Value {
	topic, _ := db.NewTopic(t.Req().Body, "i")
	topic.Author = t.Ses().Get("user").(*db.User).Id

	id, err := t.Db.Add(topic)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	// 更新余额
	conf := t.Conf()
	if conf.Reward.NewTopic != 0 {
		new(db.User).Inc(topic.Author, "balance", conf.Reward.NewTopic)
		new(db.Bill).Add(&db.Bill{
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
func (t *Topic) Edit() Value {
	topic, _ := db.NewTopic(t.Req().Body, "u")

	old := &db.Topic{}
	if err := t.Db.Find(topic.Id, old); err != nil {
		return &Fail{Msg: err.Error()}
	}

	if err := t.Db.Save(topic.Id, topic); err != nil {
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

	user := t.Ses().Get("user").(*db.User)
	if typ != 0 && old.Author != user.Id {
		new(db.Notice).Add(&db.Notice{
			Type:       uint64(typ),
			Date:       time.Now(),
			Master:     old.Author,
			User:       user.Name,
			TopicTitle: old.Title,
		})
	}

	t.Ps().Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
}

// 主题列表
func (t *Topic) List() Value {
	var s struct {
		Page int
		Node bson.ObjectId
	}
	if err := json.Unmarshal(t.Req().Body, &s); err != nil {
		return &Fail{Msg: err.Error()}
	}

	topic, err := t.Db.Paginate(s.Node, s.Page)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: topic}
}

// 主题信息
func (t *Topic) Info() Value {
	topic, _ := db.NewTopic(t.Req().Body, "b")

	if err := t.Db.Find(topic.Id, topic); err != nil {
		return &Fail{Msg: err.Error()}
	}

	t.Db.Inc(topic.Id, "views")
	return &Succ{Data: topic}
}
