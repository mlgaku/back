package module

import (
	"encoding/json"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Node struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name   string        `json:"name" validate:"required,max=30,alphanum"`
	Title  string        `json:"title" validate:"required,max=30"`
	Parent string        `json:"parent,omitempty"`
}

func (*Node) parse(body []byte) (*Node, error) {
	user := &Node{}
	return user, json.Unmarshal(body, user)
}

// 添加节点
func (n *Node) Add(ps *Pubsub, db *Database, req *Request) Value {
	node, _ := n.parse(req.Body)
	if err := com.NewVali().Struct(node); err != "" {
		return &Fail{Msg: err}
	}

	// 检查父节点
	if node.Parent != "" {
		if c, _ := db.FindId("node", node.Parent).Count(); c != 1 {
			return &Fail{Msg: "父节点不存在"}
		}
	}

	if err := db.C("node").Insert(node); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "node", Act: "list"})
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

// 获取节点信息
func (n *Node) Info(db *Database, req *Request) Value {
	node, _ := n.parse(req.Body)

	var f *mgo.Query
	if node.Id != "" {
		f = db.C("node").FindId(node.Id)
	} else if node.Name != "" {
		f = db.C("node").Find(bson.M{"name": node.Name})
	} else {
		return &Fail{Msg: "非法操作"}
	}

	if err := f.One(node); err != nil {
		return &Fail{Msg: err.Error()}
	}
	return &Succ{Data: node}
}

// 删除节点
func (n *Node) Remove(ps *Pubsub, db *Database, req *Request) Value {
	node, _ := n.parse(req.Body)

	if node.Id == "" {
		return &Fail{Msg: "ID 不能为空"}
	}

	if c, _ := db.C("node").Find(bson.M{"parent": node.Id.Hex()}).Count(); c > 0 {
		return &Fail{Msg: "删除失败: 该节点下有子节点存在"}
	}

	if err := db.C("node").RemoveId(node.Id); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "node", Act: "list"})
	return &Succ{}
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
