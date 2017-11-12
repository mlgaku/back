package db

import (
	com "github.com/mlgaku/back/common"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type Reply struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Date   time.Time     `json:"date"`
	Author bson.ObjectId `json:"author"`

	Content string        `fill:"i" json:"content" validate:"required,min=8,max=300"`
	Topic   bson.ObjectId `fill:"i" json:"topic,omitempty" validate:"required"`

	User ReplyUser `json:"user,omitempty" bson:",omitempty"`

	service.Di
}

type ReplyUser struct {
	Name   string `json:"name"`
	Avatar bool   `json:"avatar,omitempty"`
}

// 获得 Reply 实例
func NewReply(body []byte, typ string) (reply *Reply) {
	reply = &Reply{}

	if err := com.ParseJSON(body, typ, reply); err != nil {
		panic(err)
	}

	return
}

// 添加
func (r *Reply) Add(reply *Reply) {
	reply.Date = time.Now()
	reply.Content = strings.Trim(reply.Content, " ")

	if err := com.NewVali().Struct(reply); err != "" {
		panic(err)
	}

	if err := r.C("reply").Insert(reply); err != nil {
		panic(err.Error())
	}
}

// 通过作者查找
func (r *Reply) FindByAuthor(author bson.ObjectId, field M, page int) (reply []*Reply) {
	err := r.C("reply").Find(M{"author": author}).Skip(page * 20).Limit(20).Select(field).All(&reply)
	if err != nil {
		panic(err.Error())
	}

	return
}

// 分页查询
func (r *Reply) Paginate(topic bson.ObjectId, page int) (reply []*Reply) {
	if topic == "" {
		panic("主题ID不能为空")
	}

	line := []M{
		{"$match": M{"topic": topic}},
		{"$skip": page * 20},
		{"$limit": 20},
		{"$lookup": M{"from": "user", "localField": "author", "foreignField": "_id", "as": "user"}},
		{"$unwind": "$user"},
		{"$project": M{"date": 1, "content": 1, "author": 1, "user.name": 1, "user.avatar": 1}},
	}

	if err := r.C("reply").Pipe(line).All(&reply); err != nil {
		panic(err.Error())
	}

	return
}
