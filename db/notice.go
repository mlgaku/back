package db

import (
	"errors"
	. "github.com/mlgaku/back/service"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Notice struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Type   uint64        `json:"type" bson:",minsize"` // 类型(1.回复 2.At)
	Read   bool          `json:"read,omitempty"`       // 已读
	Date   time.Time     `json:"date,omitempty"`       // 日期
	Master bson.ObjectId `json:"master,omitempty"`     // 所属者ID

	Msg        string        `json:"msg,omitempty" bson:",omitempty"`                          // 通知内容
	User       string        `json:"user,omitempty" bson:",omitempty"`                         // 用户名
	TopicID    bson.ObjectId `json:"topic_id,omitempty" bson:"topic_id,omitempty"`             // (回复)主题ID
	TopicTitle string        `json:"topic_title,omitempty" bson:"topic_title,omitempty"`       // (回复)主题标题
	ReplyID    bson.ObjectId `json:"reply_id,omitempty" bson:"reply_id,omitempty"`             // (At)回复ID
	ReplyPage  uint64        `json:"reply_page,omitempty" bson:"reply_page,minsize,omitempty"` // (At)回复页数
}

// 添加
func (*Notice) Add(db *Database, notice *Notice) error {
	return db.C("notice").Insert(notice)
}

// 通过所属者查找
func (*Notice) FindByMaster(db *Database, master bson.ObjectId) (*[]Notice, error) {
	if master == "" {
		return nil, errors.New("所属者ID不能为空")
	}

	notices := &[]Notice{}
	err := db.C("notice").Find(bson.M{"read": false, "master": master}).Select(bson.M{"read": 0, "master": 0}).All(notices)
	if err != nil {
		return nil, err
	}

	return notices, nil
}

// 通过ID修改已读状态
func (*Notice) ChangeReadById(db *Database, id bson.ObjectId, read bool) error {
	if id == "" {
		return errors.New("通知ID不能为空")
	}

	return db.C("notice").UpdateId(id, bson.M{"$set": bson.M{"read": read}})
}
