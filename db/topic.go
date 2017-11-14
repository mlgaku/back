package db

import (
	com "github.com/mlgaku/back/common"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
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
	Title   string        `fill:"iu" json:"title" validate:"required,min=5,max=50"`
	Content string        `fill:"iu" json:"content,omitempty" bson:",omitempty" validate:"omitempty,min=10,max=5000"`

	User TopicUser `json:"user,omitempty" bson:",omitempty"`

	service.Di
}

type TopicUser struct {
	Name   string `json:"name,omitempty"`
	Avatar bool   `json:"avatar,omitempty"`
}

// 获得 Topic 实例
func NewTopic(body []byte, typ string) (topic *Topic) {
	topic = &Topic{}

	if err := com.ParseJSON(body, typ, topic); err != nil {
		panic(err)
	}

	return
}

// 添加
func (t *Topic) Add(topic *Topic) bson.ObjectId {
	if err := com.NewVali().Struct(topic); err != "" {
		panic(err)
	}

	if c, _ := t.C("node").FindId(topic.Node).Count(); c != 1 {
		panic("所属节点不存在")
	}

	topic.Id = bson.NewObjectId()
	topic.Date = time.Now()
	topic.Title = strings.Trim(topic.Title, " ")
	topic.Content = strings.Trim(topic.Content, " ")

	if err := t.C("topic").Insert(topic); err != nil {
		panic(err.Error())
	}

	return topic.Id
}

// 递增
func (t *Topic) Inc(id bson.ObjectId, field string) {
	if id == "" {
		panic("未指定主题ID")
	}

	if err := t.C("topic").UpdateId(id, M{"$inc": M{field: 1}}); err != nil {
		panic(err.Error())
	}
}

// 查找
func (t *Topic) Find(id bson.ObjectId) (topic *Topic) {
	if id == "" {
		panic("未指定主题ID")
	}

	if err := t.C("topic").FindId(id).One(&topic); err != nil {
		panic("主题信息获取失败")
	}

	user := new(User).Find(topic.Author, M{"name": 1, "avatar": 1})
	topic.User.Name = user.Name
	topic.User.Avatar = user.Avatar

	return
}

// 通过作者查找
func (t *Topic) FindByAuthor(author bson.ObjectId, field M, page int) (topic []*Topic) {
	err := t.C("topic").Find(M{"author": author}).Skip(page * 20).Limit(20).Select(field).All(&topic)
	if err != nil {
		panic(err.Error())
	}

	return
}

// 保存
func (t *Topic) Save(id bson.ObjectId, topic *Topic) {
	if id == "" {
		panic("主题ID不能为空")
	}

	set, err := com.Extract(topic, "u")
	if err != nil {
		panic(err.Error())
	}

	if err := t.C("topic").UpdateId(id, M{"$set": set}); err != nil {
		panic(err.Error())
	}
}

// 统计
func (t *Topic) Count(node bson.ObjectId) (c int) {
	err := error(nil)

	if node == "" {
		if c, err = t.C("topic").Count(); err != nil {
			panic(err.Error())
		}
		return
	}

	if c, err = t.C("topic").Find(M{"node": node}).Count(); err != nil {
		panic(err.Error())
	}

	return
}

// 分页查询
func (t *Topic) Paginate(node bson.ObjectId, page int, num int) (topic []*Topic) {
	if page < 1 {
		panic("页码选择不正确")
	}

	line := []M{
		{"$skip": (page - 1) * num},
		{"$limit": num},
		{"$sort": M{"date": -1}},
		{"$lookup": M{"from": "user", "localField": "author", "foreignField": "_id", "as": "user"}},
		{"$unwind": "$user"},
		{"$project": M{"date": 1, "title": 1, "node": 1, "author": 1, "views": 1, "replies": 1, "last_reply": 1, "user.name": 1, "user.avatar": 1}},
	}

	if node != "" {
		line = append([]M{
			{"$match": M{"node": node}},
		}, line[:]...)
	}

	if err := t.C("topic").Pipe(line).All(&topic); err != nil {
		panic(err.Error())
	}

	return
}

// 更新回复
func (t *Topic) UpdateReply(id bson.ObjectId, name string) {
	switch {
	case id == "":
		panic("主题ID不能为空")
	case name == "":
		panic("最后回复人名字不能为空")
	}

	err := t.C("topic").UpdateId(id, M{"$inc": M{"replies": 1}, "$set": M{"last_reply": name}})
	if err != nil {
		panic(err.Error())
	}
}
