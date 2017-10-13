package module

import (
	"encoding/json"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type Replay struct {
	Id    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Time  int64         `json:"time"`
	Topic string        `json:"topic" validate:"required,min=20,max=30"`

	Author   string `json:"author"`
	AuthorId string `json:"author_id" bson:"author_id"`

	Content string `json:"content" validate:"required,min=8,max=100"`
}

func (*Replay) parse(body []byte) (*Replay, error) {
	replay := &Replay{}
	return replay, json.Unmarshal(body, replay)
}

// 添加新回复
func (r *Replay) New(db *Database, ps *Pubsub, ses *Session, req *Request) Value {
	replay, _ := r.parse(req.Body)
	if err := com.NewVali().Struct(replay); err != "" {
		return &Fail{Msg: err}
	}

	topic := &Topic{}
	if err := db.FindId("topic", replay.Topic).One(topic); err != nil {
		return &Fail{Msg: "回复主题不存在"}
	}

	replay.Time = time.Now().Unix()
	replay.Author = ses.Get("user_name").(string)
	replay.AuthorId = ses.Get("user_id").(bson.ObjectId).Hex()
	replay.Content = strings.Trim(replay.Content, " ")

	if err := db.C("replay").Insert(replay); err != nil {
		return &Fail{Msg: err.Error()}
	}

	// 回复人不是主题作者时添加通知
	if replay.AuthorId != topic.AuthorId {
		db.C("notice").Insert(&Notice{
			Type:       1,
			Time:       time.Now().Unix(),
			Master:     topic.AuthorId,
			User:       replay.Author,
			TopicID:    replay.Topic,
			TopicTitle: topic.Title,
		})
	}

	ps.Publish(&Prot{Mod: "replay", Act: "list"})
	ps.Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
}

// 获取回复列表
func (r *Replay) List(db *Database, req *Request) Value {
	var s struct {
		Page  int
		Topic string
	}
	if err := json.Unmarshal(req.Body, &s); err != nil {
		return &Fail{Msg: err.Error()}
	}

	if s.Topic == "" {
		return &Fail{Msg: "主题ID不能为空"}
	}

	replay := &[]Replay{}
	db.C("replay").Find(bson.M{"topic": s.Topic}).Skip(s.Page * 20).Limit(20).All(replay)
	return &Succ{Data: replay}
}
