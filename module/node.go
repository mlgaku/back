package module

import (
	"encoding/json"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
)

type Node struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name   string        `json:"name" validate:"required,max=30,alphanum"`
	Title  string        `json:"title" validate:"required,max=30"`
	Parent string        `json:"parent,omitempty" validate:"omitempty,max=30,alphanum"`
}

func (*Node) parse(body []byte) (*Node, error) {
	user := &Node{}
	return user, json.Unmarshal(body, user)
}

// 添加节点
func (n *Node) Add(db *Database, req *Request) Value {
	node, _ := n.parse(req.Body)
	if err := com.NewVali().Struct(node); err != "" {
		return &Fail{Msg: err}
	}

	// 检查父节点
	if node.Parent != "" {
		if c, _ := db.C("node").Find(bson.M{"name": node.Parent}).Count(); c < 1 {
			return &Fail{Msg: "父节点不存在"}
		}
	}

	if err := db.C("node").Insert(node); err != nil {
		return &Fail{Msg: err.Error()}
	}
	return &Succ{}
}

// 获取节点列表
func (n *Node) List(db *Database) Value {
	node := &[]Node{}
	if err := db.C("node").Find(bson.M{}).All(node); err != nil {
		return &Fail{Msg: err.Error()}
	}
	return &Succ{Data: node}
}

// 检查是否有相同节点存在
func (n *Node) Check(db *Database, req *Request) Value {
	node, _ := n.parse(req.Body)

	if err := com.NewVali().Var(node.Name, "required"); err != "" {
		return &Fail{Msg: err}
	}

	if c, _ := db.C("node").Find(bson.M{"name": node.Name}).Count(); c > 0 {
		return &Fail{Msg: "已有同名节点存在"}
	}

	return &Succ{}
}
