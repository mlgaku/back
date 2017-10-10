package module

import (
	"fmt"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
)

type Topic struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`

	Time    int64  `json:"time"`
	Title   string `json:"title"`
	Content string `json:"content,omitempty"`

	Author   string `json:"author"`
	AuthorId string `json:"author_id"`

	Views   uint64 `json:"views"`
	Replies uint64 `json:"replies"`
}

func (*Topic) parse() {

}

// 发表新主题
func (t *Topic) New() {

}

func (t *Topic) Test1(ses *Session) {
	ses.Set("name", "yazi")
}

func (t *Topic) Test2(ses *Session) Value {
	fmt.Println(ses.Get("name"))
	return ses.Get("name")
}
