package module

import (
	"encoding/json"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
)

type Notice struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Type   uint64        `json:"type" bson:",minsize"` // 类型(1.回复 2.At)
	Read   bool          `json:"read,omitempty"`       // 已读
	Time   int64         `json:"time,omitempty"`       // 时间
	Master bson.ObjectId `json:"master,omitempty"`     // 所属者ID

	Msg        string        `json:"msg,omitempty" bson:",omitempty"`                            // 通知内容
	User       string        `json:"user,omitempty" bson:",omitempty"`                           // 用户名
	TopicID    bson.ObjectId `json:"topic_id,omitempty" bson:"topic_id,omitempty"`               // (回复)主题ID
	TopicTitle string        `json:"topic_title,omitempty" bson:"topic_title,omitempty"`         // (回复)主题标题
	ReplayID   bson.ObjectId `json:"replay_id,omitempty" bson:"replay_id,omitempty"`             // (At)回复ID
	ReplayPage uint64        `json:"replay_page,omitempty" bson:"replay_page,minsize,omitempty"` // (At)回复页数
}

func (*Notice) parse(body []byte) (*Notice, error) {
	replay := &Notice{}
	return replay, json.Unmarshal(body, replay)
}

// 获取通知列表
func (n *Notice) List(db *Database, req *Request) Value {
	notice, _ := n.parse(req.Body)
	if notice.Master == "" {
		return &Fail{Msg: "所属者ID不能为空"}
	}

	notices := &[]Notice{}
	db.C("notice").Find(bson.M{"read": false, "master": notice.Master}).Select(bson.M{"read": 0, "master": 0}).All(notices)
	return &Succ{Data: notices}
}

// 移除通知
func (n *Notice) Remove(ps *Pubsub, db *Database, req *Request) Value {
	notice, _ := n.parse(req.Body)
	if notice.Id == "" {
		return &Fail{Msg: "通知ID不能为空"}
	}

	if err := db.C("notice").UpdateId(notice.Id, bson.M{"$set": bson.M{"read": true}}); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
}
