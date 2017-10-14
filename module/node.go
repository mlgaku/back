package module

import (
	"encoding/json"
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Node struct {
	Db db.Node
}

func (*Node) parse(body []byte) (*db.Node, error) {
	user := &db.Node{}
	return user, json.Unmarshal(body, user)
}

// 添加节点
func (n *Node) Add(ps *Pubsub, db *Database, req *Request) Value {
	node, _ := n.parse(req.Body)

	// 检查父节点
	if node.Parent != "" {
		if b, err := n.Db.IdExists(db, node.Parent); err != nil {
			return &Fail{Msg: err.Error()}
		} else if !b {
			return &Fail{Msg: "父节点不存在"}
		}
	}

	if err := n.Db.Add(db, node); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "node", Act: "list"})
	return &Succ{}
}

// 获取节点列表
func (n *Node) List(db *Database) Value {
	node, err := n.Db.FindAll(db)
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: node}
}

// 获取节点信息
func (n *Node) Info(db *Database, req *Request) Value {
	node, _ := n.parse(req.Body)
	if err := n.Db.FindByIdOrName(db, node); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: node}
}

// 删除节点
func (n *Node) Remove(ps *Pubsub, db *Database, req *Request) Value {
	node, _ := n.parse(req.Body)

	// 检查子节点
	if b, err := n.Db.HasChild(db, node.Id); err != nil {
		return &Fail{Msg: err.Error()}
	} else if b {
		return &Fail{Msg: "删除失败: 该节点下有子节点存在"}
	}

	if err := n.Db.RemoveById(db, node.Id); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "node", Act: "list"})
	return &Succ{}
}

// 检查是否有相同节点存在
func (n *Node) Check(db *Database, req *Request) Value {
	node, _ := n.parse(req.Body)
	if b, err := n.Db.NameExists(db, node.Name); err != nil {
		return &Fail{Msg: err.Error()}
	} else if b {
		return &Fail{Msg: "已有同名节点存在"}
	}

	return &Succ{}
}
