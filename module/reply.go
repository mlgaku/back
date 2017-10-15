package module

import (
	"encoding/json"
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
	"regexp"
	"strings"
	"time"
)

type Reply struct {
	Db db.Reply
}

func (*Reply) parse(body []byte) (*db.Reply, error) {
	reply := &db.Reply{}
	return reply, json.Unmarshal(body, reply)
}

// 添加新回复
func (r *Reply) New(bd *Database, ps *Pubsub, ses *Session, req *Request) Value {
	reply, _ := r.parse(req.Body)
	reply.Author = ses.Get("user_id").(bson.ObjectId)

	topic := &db.Topic{}
	if err := topic.Find(bd, reply.Topic, topic); err != nil {
		return &Fail{Msg: err.Error()}
	}

	// 添加回复
	if err := r.Db.Add(bd, reply); err != nil {
		return &Fail{Msg: err.Error()}
	}

	// 更新最后回复
	topic.UpdateReply(bd, reply.Topic, ses.Get("user_name").(string))

	// 回复人不是主题作者时添加通知
	if reply.Author != topic.Author {
		new(db.Notice).Add(bd, &db.Notice{
			Type:       1,
			Date:       time.Now(),
			Master:     topic.Author,
			User:       ses.Get("user_name").(string),
			TopicID:    reply.Topic,
			TopicTitle: topic.Title,
		})
	}

	// 通知被at的人
	r.handleAt(bd, ses, topic, reply)

	ps.Publish(&Prot{Mod: "reply", Act: "list"})
	ps.Publish(&Prot{Mod: "notice", Act: "list"})
	return &Succ{}
}

// 获取回复列表
func (r *Reply) List(db *Database, req *Request) Value {
	var s struct {
		Page  int
		Topic bson.ObjectId
	}
	if err := json.Unmarshal(req.Body, &s); err != nil {
		return &Fail{Msg: err.Error()}
	}

	reply, err := r.Db.Paginate(db, s.Topic, s.Page)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: reply}
}

// 处理 At
func (*Reply) handleAt(bd *Database, ses *Session, topic *db.Topic, reply *db.Reply) {
	match := regexp.MustCompile(`@[a-zA-Z0-9]+`).FindAllString(reply.Content, 5)
	if match == nil {
		return
	}

	for k, v := range match {
		match[k] = strings.TrimLeft(v, "@")
	}

	user, err := new(db.User).FindByNameMany(bd, match)
	if err != nil {
		return
	}

	name := ses.Get("user_name").(string)
	for k, v := range user {
		// 跳过@自己
		if k == name {
			continue
		}

		new(db.Notice).Add(bd, &db.Notice{
			Type:       2,
			Date:       time.Now(),
			Master:     v.Id,
			User:       name,
			TopicID:    reply.Topic,
			TopicTitle: topic.Title,
		})
	}
}
