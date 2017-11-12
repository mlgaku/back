package module

import (
	"github.com/mlgaku/back/db"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
)

type Node struct {
	db db.Node

	service.Di
}

// 添加节点
func (n *Node) Add() Value {
	node := db.NewNode(n.Req().Body, "i")

	// 检查父节点
	if node.Parent != "" {
		if b := n.db.IdExists(node.Parent); !b {
			return &Fail{Msg: "父节点不存在"}
		}
	}

	n.db.Add(node)

	n.Ps().Publish(&Prot{Mod: "node", Act: "list"})
	return &Succ{}
}

// 编辑节点
func (n *Node) Edit() Value {
	node := db.NewNode(n.Req().Body, "u")

	// 检查父节点
	if node.Parent != "" {
		if b := n.db.IdExists(node.Parent); !b {
			return &Fail{Msg: "父节点不存在"}
		}
	}

	n.db.Save(node.Id, node)

	n.Ps().Publish(&Prot{Mod: "node", Act: "list"})
	return &Succ{}
}

// 获取节点列表
func (n *Node) List() Value {
	return &Succ{Data: n.db.FindAll()}
}

// 获取节点信息
func (n *Node) Info() Value {
	node := db.NewNode(n.Req().Body, "b")
	n.db.FindByIdOrName(node)

	return &Succ{Data: node}
}

// 删除节点
func (n *Node) Remove() Value {
	node := db.NewNode(n.Req().Body, "b")

	// 检查子节点
	if b := n.db.HasChild(node.Id); b {
		return &Fail{Msg: "删除失败：该节点下有子节点存在"}
	}

	// 检查节点下是否还有主题存在
	if b := n.db.HasTopic(node.Id); b {
		return &Fail{Msg: "删除失败：该节点下还有主题存在"}
	}

	n.db.RemoveById(node.Id)

	n.Ps().Publish(&Prot{Mod: "node", Act: "list"})
	return &Succ{}
}

// 检查节点名是否可用
func (n *Node) Check() Value {
	node := db.NewNode(n.Req().Body, "b")
	return &Succ{Data: !n.db.NameExists(node.Name)}
}
