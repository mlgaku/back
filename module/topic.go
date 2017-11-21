package module

import (
	"encoding/json"
	"github.com/mlgaku/back/db"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
	"math"
	"time"
)

type Topic struct {
	db db.Topic

	service.Di
}

// 发表新主题
func (t *Topic) New() Value {
	topic := db.NewTopic(t.Req().Body, "i")
	topic.Author = t.Ses().Get("user").(*db.User).Id

	id := t.db.Add(topic)

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
	topic := db.NewTopic(t.Req().Body, "u")

	old := t.db.Find(topic.Id)
	t.db.Save(topic.Id, topic)

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

	return &Succ{Data: M{
		"per":   20,
		"page":  s.Page,
		"total": math.Ceil(float64(t.db.Count(s.Node)) / 20),
		"list":  t.db.Paginate(s.Node, s.Page, 20),
	}}
}

// 主题信息
func (t *Topic) Info() Value {
	topic := t.db.Find(db.NewTopic(t.Req().Body, "b").Id)

	t.db.Inc(topic.Id, "views")
	return &Succ{Data: topic}
}

// 补充内容
func (t *Topic) Subtle() Value {
	m := map[string]string{}
	if err := json.Unmarshal(t.Di.Req().Body, &m); err != nil {
		return &Fail{Msg: err.Error()}
	}

	topic := bson.ObjectIdHex(m["topic"])
	if t.db.Find(topic).Author != t.Di.Ses().Get("user").(*db.User).Id {
		return &Fail{Msg: "你不能给别人的主题补充内容"}
	}

	t.db.AddSubtle(topic, &db.TopicSubtle{Content: m["content"]})
	return &Succ{}
}
