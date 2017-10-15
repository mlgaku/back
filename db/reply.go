package db

import (
	"errors"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type Reply struct {
	Id      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Time    time.Time     `json:"time"`
	Content string        `json:"content" validate:"required,min=8,max=300"`

	Topic  bson.ObjectId `json:"topic,omitempty" validate:"required"`
	Author bson.ObjectId `json:"author"`

	User struct {
		Name string `json:"name"`
	} `json:"user,omitempty" bson:",omitempty"`
}

// 添加
func (*Reply) Add(db *Database, reply *Reply) error {
	if err := com.NewVali().Struct(reply); err != "" {
		return errors.New(err)
	}

	reply.Time = time.Now()
	reply.Content = strings.Trim(reply.Content, " ")

	return db.C("reply").Insert(reply)
}

// 分页查询
func (*Reply) Paginate(db *Database, topic bson.ObjectId, page int) (*[]Reply, error) {
	if topic == "" {
		return nil, errors.New("主题ID不能为空")
	}

	line := []bson.M{
		bson.M{"$match": bson.M{"topic": topic}},
		bson.M{"$skip": page * 20},
		bson.M{"$limit": 20},
		bson.M{"$lookup": bson.M{"from": "user", "localField": "author", "foreignField": "_id", "as": "user"}},
		bson.M{"$unwind": "$user"},
		bson.M{"$project": bson.M{"time": 1, "content": 1, "author": 1, "user.name": 1}},
	}

	reply := &[]Reply{}
	if err := db.C("reply").Pipe(line).All(reply); err != nil {
		return nil, err
	}

	return reply, nil
}
