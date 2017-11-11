package module

import (
	"github.com/mlgaku/back/db"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Node struct {
	Db db.Node
	service.Di
}

// 添加节点
func (n *Node) Add() Value {
	node, _ := db.NewNode(n.Req().Body, "i")

	// 检查父节点
	if node.Parent != "" {
		if b, err := n.Db.IdExists(node.Parent); err != nil {
			return &Fail{Msg: err.Error()}
		} else if !b {
			return &Fail{Msg: "父节点不存在"}
		}
	}

	if err := n.Db.Add(node); err != nil {
		return &Fail{Msg: err.Error()}
	}

	n.Ps().Publish(&Prot{Mod: "node", Act: "list"})
	return &Succ{}
}

// 编辑节点
func (n *Node) Edit() Value {
	node, _ := db.NewNode(n.Req().Body, "u")

	// 检查父节点
	if node.Parent != "" {
		if b, err := n.Db.IdExists(node.Parent); err != nil {
			return &Fail{Msg: err.Error()}
		} else if !b {
			return &Fail{Msg: "父节点不存在"}
		}
	}

	if err := n.Db.Save(node.Id, node); err != nil {
		return &Fail{Msg: err.Error()}
	}

	n.Ps().Publish(&Prot{Mod: "node", Act: "list"})
	return &Succ{}
}

// 获取节点列表
func (n *Node) List() Value {
	node, err := n.Db.FindAll()
	if err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: node}
}

// 获取节点信息
func (n *Node) Info() Value {
	node, _ := db.NewNode(n.Req().Body, "b")
	if err := n.Db.FindByIdOrName(node); err != nil {
		return &Fail{Msg: err.Error()}
	}

	return &Succ{Data: node}
}

// 删除节点
func (n *Node) Remove() Value {
	node, _ := db.NewNode(n.Req().Body, "b")

	// 检查子节点
	if b, err := n.Db.HasChild(node.Id); err != nil {
		return &Fail{Msg: err.Error()}
	} else if b {
		return &Fail{Msg: "删除失败: 该节点下有子节点存在"}
	}

	// 检查节点下是否还有主题存在
	if b, err := n.Db.HasTopic(node.Id); err != nil {
		return &Fail{Msg: err.Error()}
	} else if b {
		return &Fail{Msg: "删除失败: 该节点下还有主题存在"}
	}

	if err := n.Db.RemoveById(node.Id); err != nil {
		return &Fail{Msg: err.Error()}
	}

	n.Ps().Publish(&Prot{Mod: "node", Act: "list"})
	return &Succ{}
}

// 检查节点名是否可用
func (n *Node) Check() Value {
	node, _ := db.NewNode(n.Req().Body, "b")
	if b, err := n.Db.NameExists(node.Name); err != nil {
		return &Fail{Msg: err.Error()}
	} else {
		return &Succ{Data: !b}
	}
}
