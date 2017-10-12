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
	Read   bool          `json:"read"`                 // 已读
	Master string        `json:"master,omitempty"`     // 所属者ID

	Msg        string `json:"msg,omitempty" bson:",omitempty"`                            // 通知内容
	User       string `json:"user,omitempty" bson:",omitempty"`                           // 用户名
	TopicID    string `json:"topic_id,omitempty" bson:"topic_id,omitempty"`               // (回复)主题ID
	TopicTitle string `json:"topic_title,omitempty" bson:"topic_title,omitempty"`         // (回复)主题标题
	ReplayID   string `json:"replay_id,omitempty" bson:"replay_id,omitempty"`             // (At)回复ID
	ReplayPage uint64 `json:"replay_page,omitempty" bson:"replay_page,minsize,omitempty"` // (At)回复页数
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
	db.C("notice").Find(bson.M{"master": notice.Master}).Select(bson.M{"master": 0}).All(notices)
	return &Succ{Data: notices}
}
