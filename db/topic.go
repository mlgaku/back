package db

import (
	"errors"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type Topic struct {
	Date      time.Time     `json:"date"`
	Author    bson.ObjectId `json:"author"`
	Views     uint64        `json:"views"`
	Replies   uint64        `json:"replies"`
	LastReply string        `json:"last_reply,omitempty" bson:"last_reply,omitempty"`

	Id      bson.ObjectId `fill:"u" json:"id" bson:"_id,omitempty"`
	Node    bson.ObjectId `fill:"iu" json:"node" validate:"required"`
	Title   string        `fill:"iu" json:"title" validate:"required,min=10,max=50"`
	Content string        `fill:"iu" json:"content,omitempty" bson:",omitempty" validate:"omitempty,min=20,max=5000"`

	User TopicUser `json:"user,omitempty" bson:",omitempty"`
}

type TopicUser struct {
	Name   string `json:"name,omitempty"`
	Avatar bool   `json:"avatar,omitempty"`
}

// 获得 Topic 实例
func NewTopic(body []byte, typ string) (*Topic, error) {
	topic := &Topic{}
	if err := com.ParseJSON(body, typ, topic); err != nil {
		panic(err)
	}

	return topic, nil
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
	topic.Date = time.Now()
	topic.Title = strings.Trim(topic.Title, " ")
	topic.Content = strings.Trim(topic.Content, " ")

	if err := db.C("topic").Insert(topic); err != nil {
		return "", err
	}

	return topic.Id, nil
}

// 递增
func (*Topic) Inc(db *Database, id bson.ObjectId, field string) error {
	if id == "" {
		return errors.New("未指定主题ID")
	}
	return db.C("topic").UpdateId(id, bson.M{"$inc": bson.M{field: 1}})
}

// 查找
func (*Topic) Find(db *Database, id bson.ObjectId, topic *Topic) error {
	if id == "" {
		return errors.New("未指定主题ID")
	}

	if err := db.C("topic").FindId(id).One(topic); err != nil {
		return errors.New("主题信息获取失败")
	}

	if err := new(User).Find(db, topic.Author, &topic.User, bson.M{
		"name":   1,
		"avatar": 1,
	}); err != nil {
		return errors.New("用户信息获取失败")
	}

	return nil
}

// 通过作者查找
func (*Topic) FindByAuthor(db *Database, author bson.ObjectId, field bson.M, page int) (*[]Topic, error) {
	result := new([]Topic)
	return result, db.C("topic").Find(bson.M{"author": author}).Skip(page * 20).Limit(20).Select(field).All(result)
}

// 保存
func (*Topic) Save(db *Database, id bson.ObjectId, topic *Topic) error {
	if id == "" {
		return errors.New("主题ID不能为空")
	}

	set, err := com.Extract(topic, "u")
	if err != nil {
		return err
	}

	return db.C("topic").UpdateId(id, bson.M{"$set": set})
}

// 分页查询
func (*Topic) Paginate(db *Database, node bson.ObjectId, page int) (*[]Topic, error) {
	line := []bson.M{
		{"$skip": page * 20},
		{"$limit": 20},
		{"$lookup": bson.M{"from": "user", "localField": "author", "foreignField": "_id", "as": "user"}},
		{"$unwind": "$user"},
		{"$project": bson.M{"date": 1, "title": 1, "node": 1, "author": 1, "views": 1, "replies": 1, "last_reply": 1, "user.name": 1, "user.avatar": 1}},
	}

	if node != "" {
		line = append([]bson.M{
			{"$match": bson.M{"node": node}},
		}, line[:]...)
	}

	topic := &[]Topic{}
	if err := db.C("topic").Pipe(line).All(topic); err != nil {
		return nil, err
	}

	return topic, nil
}

// 更新回复
func (*Topic) UpdateReply(db *Database, id bson.ObjectId, name string) error {
	switch {
	case id == "":
		return errors.New("主题ID不能为空")
	case name == "":
		return errors.New("最后回复人名字不能为空")
	}

	return db.C("topic").UpdateId(id, bson.M{"$inc": bson.M{"replies": 1}, "$set": bson.M{"last_reply": name}})
}
