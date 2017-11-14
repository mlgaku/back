package db

import (
	com "github.com/mlgaku/back/common"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Notice struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Type   uint64        `json:"type" bson:",minsize"` // 类型(1.回复 2.At 3.修改主题 4.移动主题 5.修改&移动主题)
	Read   bool          `json:"read,omitempty"`       // 已读
	Date   time.Time     `json:"date,omitempty"`       // 日期
	Master bson.ObjectId `json:"master,omitempty"`     // 所属者ID

	Msg        string        `json:"msg,omitempty" bson:",omitempty"`                          // 通知内容
	User       string        `json:"user,omitempty" bson:",omitempty"`                         // 用户名
	TopicID    bson.ObjectId `json:"topic_id,omitempty" bson:"topic_id,omitempty"`             // (回复)主题ID
	TopicTitle string        `json:"topic_title,omitempty" bson:"topic_title,omitempty"`       // (回复)主题标题
	ReplyID    bson.ObjectId `json:"reply_id,omitempty" bson:"reply_id,omitempty"`             // (At)回复ID
	ReplyPage  uint64        `json:"reply_page,omitempty" bson:"reply_page,minsize,omitempty"` // (At)回复页数

	service.Di `json:"-" bson:"-"`
}

// 获得 Notice 实例
func NewNotice(body []byte, typ string) (notice *Notice) {
	notice = &Notice{}

	if err := com.ParseJSON(body, typ, notice); err != nil {
		panic(err)
	}

	return
}

// 添加
func (n *Notice) Add(notice *Notice) {
	if err := n.C("notice").Insert(notice); err != nil {
		panic(err.Error())
	}
}

// 查找
func (n *Notice) Find(id bson.ObjectId) (notice *Notice) {
	if id == "" {
		panic("未指定通知ID")
	}

	if err := n.C("notice").FindId(id).One(&notice); err != nil {
		panic(err.Error())
	}

	return
}

// 通过所属者查找
func (n *Notice) FindByMaster(master bson.ObjectId) (notices []*Notice) {
	if master == "" {
		panic("所属者ID不能为空")
	}

	err := n.C("notice").Find(M{"read": false, "master": master}).Select(M{"read": 0, "master": 0}).All(&notices)
	if err != nil {
		panic(err.Error())
	}

	return
}

// 通过ID修改已读状态
func (n *Notice) ChangeReadById(id bson.ObjectId, read bool) {
	if id == "" {
		panic("通知ID不能为空")
	}

	if err := n.C("notice").UpdateId(id, M{"$set": M{"read": read}}); err != nil {
		panic(err.Error())
	}
}
