package module

import (
	"encoding/json"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type Topic struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`

	Node    string `json:"node" validate:"required,min=20,max=30"`
	Time    int64  `json:"time"`
	Title   string `json:"title" validate:"required,min=10,max=50"`
	Content string `json:"content,omitempty" bson:",omitempty" validate:"omitempty,required,min=30,max=1000"`

	Author   string `json:"author"`
	AuthorId string `json:"author_id"`

	Views   uint64 `json:"views"`
	Replies uint64 `json:"replies"`
}

func (*Topic) parse(body []byte) (*Topic, error) {
	topic := &Topic{}
	return topic, json.Unmarshal(body, topic)
}

// 发表新主题
func (t *Topic) New(db *Database, ses *Session, req *Request) Value {
	topic, _ := t.parse(req.Body)

	if err := com.NewVali().Struct(topic); err != "" {
		return &Fail{Msg: err}
	}

	if c, _ := db.FindId("node", topic.Node).Count(); c != 1 {
		return &Fail{Msg: "所选节点不存在"}
	}

	topic.Id = bson.NewObjectId()
	topic.Time = time.Now().Unix()
	topic.Title = strings.Trim(topic.Title, " ")
	topic.Content = strings.Trim(topic.Content, " ")
	topic.Author = ses.Get("user_name").(string)
	topic.AuthorId = ses.Get("user_id").(bson.ObjectId).Hex()

	if err := db.C("topic").Insert(topic); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: topic.Id}
}

// 主题列表
func (t *Topic) List(db *Database, req *Request) Value {
	var s struct {
		Page int
		Node string
	}
	if err := json.Unmarshal(req.Body, &s); err != nil {
		return &Fail{Msg: err.Error()}
	}

	var q *mgo.Query
	if s.Node == "" {
		q = db.C("topic").Find(nil)
	} else {
		q = db.C("topic").Find(bson.M{"node": s.Node})
	}

	topic := &[]Topic{}
	q.Skip(s.Page * 20).Limit(20).All(topic)
	return &Succ{Data: topic}
}

// 主题信息
func (t *Topic) Info(db *Database, req *Request) Value {
	topic, _ := t.parse(req.Body)
	if topic.Id == "" {
		return &Fail{Msg: "未指定主题ID"}
	}

	if err := db.C("topic").FindId(topic.Id).One(topic); err != nil {
		return &Fail{Msg: "主题信息获取失败"}
	}

	return &Succ{Data: topic}
}
