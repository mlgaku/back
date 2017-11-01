package module

import (
	"github.com/mlgaku/back/db"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Node struct {
	Db db.Node
}

// 添加节点
func (n *Node) Add(ps *Pubsub, bd *Database, req *Request) Value {
	node, _ := db.NewNode(req.Body, "i")

	// 检查父节点
	if node.Parent != "" {
		if b, err := n.Db.IdExists(bd, node.Parent); err != nil {
			return &Fail{Msg: err.Error()}
		} else if !b {
			return &Fail{Msg: "父节点不存在"}
		}
	}

	if err := n.Db.Add(bd, node); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "node", Act: "list"})
	return &Succ{}
}

// 编辑节点
func (n *Node) Edit(ps *Pubsub, bd *Database, req *Request) Value {
	node, _ := db.NewNode(req.Body, "u")

	// 检查父节点
	if node.Parent != "" {
		if b, err := n.Db.IdExists(bd, node.Parent); err != nil {
			return &Fail{Msg: err.Error()}
		} else if !b {
			return &Fail{Msg: "父节点不存在"}
		}
	}

	if err := n.Db.Save(bd, node.Id, node); err != nil {
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
func (n *Node) Info(bd *Database, req *Request) Value {
	node, _ := db.NewNode(req.Body, "b")
	if err := n.Db.FindByIdOrName(bd, node); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: node}
}

// 删除节点
func (n *Node) Remove(ps *Pubsub, bd *Database, req *Request) Value {
	node, _ := db.NewNode(req.Body, "b")

	// 检查子节点
	if b, err := n.Db.HasChild(bd, node.Id); err != nil {
		return &Fail{Msg: err.Error()}
	} else if b {
		return &Fail{Msg: "删除失败: 该节点下有子节点存在"}
	}

	if err := n.Db.RemoveById(bd, node.Id); err != nil {
		return &Fail{Msg: err.Error()}
	}

	ps.Publish(&Prot{Mod: "node", Act: "list"})
	return &Succ{}
}

// 检查节点名是否可用
func (n *Node) Check(bd *Database, req *Request) Value {
	node, _ := db.NewNode(req.Body, "b")
	if b, err := n.Db.NameExists(bd, node.Name); err != nil {
		return &Fail{Msg: err.Error()}
	} else {
		return &Succ{Data: !b}
	}
}
