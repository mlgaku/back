package db

import (
	"errors"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type Replay struct {
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
func (*Replay) Add(db *Database, replay *Replay) error {
	if err := com.NewVali().Struct(replay); err != "" {
		return errors.New(err)
	}

	replay.Time = time.Now()
	replay.Content = strings.Trim(replay.Content, " ")
	if err := db.C("replay").Insert(replay); err != nil {
		return err
	}

	return nil
}

// 分页查询
func (*Replay) Paginate(db *Database, topic bson.ObjectId, page int) (*[]Replay, error) {
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

	replay := &[]Replay{}
	if err := db.C("replay").Pipe(line).All(replay); err != nil {
		return nil, err
	}

	return replay, nil
}
