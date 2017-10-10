package module

import (
	"encoding/json"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Topic struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`

	Node    string `json:"node" validate:"required,min=20,max=30"`
	Time    int64  `json:"time"`
	Title   string `json:"title" validate:"required,min=10,max=50"`
	Content string `json:"content,omitempty" validate:"omitempty,required,min=30,max=1000"`

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
	topic.Author = ses.Get("user_name").(string)
	topic.AuthorId = ses.Get("user_id").(bson.ObjectId).Hex()

	if err := db.C("topic").Insert(topic); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: topic.Id}
}
