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

	Topic      bson.ObjectId `json:"topic" validate:"required"`
	Author     bson.ObjectId `json:"author"`
	AuthorName string        `json:"author_name" bson:"author_name"`
}

// 添加
func (*Replay) Add(db *Database, replay *Replay) error {
	if err := com.NewVali().Struct(replay); err != "" {
		return errors.New(err)
	}

	topic := &Topic{}
	if err := db.C("topic").FindId(replay.Topic).One(topic); err != nil {
		return errors.New("回复的主题不存在")
	}

	replay.Time = time.Now()
	replay.Content = strings.Trim(replay.Content, " ")
	if err := db.C("replay").Insert(replay); err != nil {
		return err
	}

	// 回复人不是主题作者时添加通知
	if replay.Author != topic.Author {
		new(Notice).Add(db, &Notice{
			Type:       1,
			Time:       time.Now(),
			Master:     topic.Author,
			User:       replay.AuthorName,
			TopicID:    replay.Topic,
			TopicTitle: topic.Title,
		})

	}

	return nil
}

// 分页查询
func (*Replay) Paginate(db *Database, topic bson.ObjectId, page int) (*[]Replay, error) {
	if topic == "" {
		return nil, errors.New("主题ID不能为空")
	}

	replay := &[]Replay{}
	if err := db.C("replay").Find(bson.M{"topic": topic}).Skip(page * 20).Limit(20).All(replay); err != nil {
		return nil, err
	}

	return replay, nil
}
