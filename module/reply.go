package module

import (
	"encoding/json"
	"github.com/mlgaku/back/db"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"strings"
	"time"
)

type Reply struct {
	db db.Reply
	service.Di
}

// 添加新回复
func (r *Reply) New() Value {
	user := r.Ses().Get("user").(*db.User)

	reply, _ := db.NewReply(r.Req().Body, "i")
	reply.Author = user.Id

	topic := &db.Topic{}
	if err := topic.Find(reply.Topic, topic); err != nil {
		return &Fail{Msg: err.Error()}
	}

	// 添加回复
	if err := r.db.Add(reply); err != nil {
		return &Fail{Msg: err.Error()}
	}

	// 更新最后回复
	topic.UpdateReply(reply.Topic, user.Name)

	// 回复人不是主题作者时添加通知
	if reply.Author != topic.Author {
		new(db.Notice).Add(&db.Notice{
			Type:       1,
			Date:       time.Now(),
			Master:     topic.Author,
			User:       user.Name,
			TopicID:    reply.Topic,
			TopicTitle: topic.Title,
		})
	}

	// 通知被at的人
	r.handleAt(user.Name, topic, reply)

	// 更新余额
	conf := r.Conf()
	if conf.Reward.NewReply != 0 {
		new(db.User).Inc(reply.Author, "balance", conf.Reward.NewReply)
		new(db.Bill).Add(&db.Bill{
			Msg:    topic.Title,
			Type:   2,
			Date:   time.Now(),
			Number: conf.Reward.NewReply,
			Master: reply.Author,
		})
	}

	r.Ps().Publish(&Prot{Mod: "reply", Act: "list"})
	r.Ps().Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
}

// 获取回复列表
func (r *Reply) List() Value {
	var s struct {
		Page  int
		Topic bson.ObjectId
	}
	if err := json.Unmarshal(r.Req().Body, &s); err != nil {
		return &Fail{Msg: err.Error()}
	}

	reply, err := r.db.Paginate(s.Topic, s.Page)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: reply}
}

// 处理 At
func (*Reply) handleAt(name string, topic *db.Topic, reply *db.Reply) {
	match := regexp.MustCompile(`@[a-zA-Z0-9]+`).FindAllString(reply.Content, 5)
	if match == nil {
		return
	}

	for k, v := range match {
		match[k] = strings.TrimLeft(v, "@")
	}

	user, err := new(db.User).FindByNameMany(match)
	if err != nil {
		return
	}

	for k, v := range user {
		// 跳过@自己
		if k == name {
			continue
		}

		new(db.Notice).Add(&db.Notice{
			Type:       2,
			Date:       time.Now(),
			Master:     v.Id,
			User:       name,
			TopicID:    reply.Topic,
			TopicTitle: topic.Title,
		})
	}
}
