package db

import (
	"errors"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type Topic struct {
	Id      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Time    time.Time     `json:"time"`
	Title   string        `json:"title" validate:"required,min=10,max=50"`
	Content string        `json:"content,omitempty" bson:",omitempty" validate:"omitempty,required,min=20,max=5000"`

	Node   bson.ObjectId `json:"node" validate:"required"`
	Author bson.ObjectId `json:"author"`

	Views   uint64 `json:"views"`
	Replies uint64 `json:"replies"`
}

// 添加
func (*Topic) Add(db *Database, topic *Topic) (bson.ObjectId, error) {

	if err := com.NewVali().Struct(topic); err != "" {
		return "", errors.New(err)
	}

	if c, _ := db.C("node").FindId(topic.Node).Count(); c != 1 {
		return "", errors.New("所属节点不存在")
	}

	topic.Id = bson.NewObjectId()
	topic.Time = time.Now()
	topic.Title = strings.Trim(topic.Title, " ")
	topic.Content = strings.Trim(topic.Content, " ")

	if err := db.C("topic").Insert(topic); err != nil {
		return "", err
	}

	return topic.Id, nil
}

// 查询
func (*Topic) Find(db *Database, id bson.ObjectId, topic *Topic) error {
	if id == "" {
		return errors.New("未指定主题ID")
	}

	if err := db.C("topic").FindId(id).One(topic); err != nil {
		return errors.New("主题信息获取失败")
	}

	return nil
}

// 分页查询
func (*Topic) Paginate(db *Database, node bson.ObjectId, page int) (*[]Topic, error) {
	var q *mgo.Query
	if node == "" {
		q = db.C("topic").Find(nil)
	} else {
		q = db.C("topic").Find(bson.M{"node": node})
	}

	topic := &[]Topic{}
	if err := q.Skip(page * 20).Limit(20).Select(bson.M{"content": 0}).All(topic); err != nil {
		return nil, err
	}

	return topic, nil
}
