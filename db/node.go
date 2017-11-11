package db

import (
	"errors"
	com "github.com/mlgaku/back/common"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Node struct {
	Id     bson.ObjectId `fill:"u" json:"id" bson:"_id,omitempty"`
	Name   string        `fill:"iu" json:"name" validate:"required,max=30,alphanum"`                            // 名字
	Title  string        `fill:"iu" json:"title" validate:"required,max=30"`                                    // 标题
	Sort   uint64        `fill:"iu" json:"sort,omitempty" bson:",omitempty"`                                    // 排序
	Desc   string        `fill:"iu" json:"desc,omitempty" bson:",omitempty" validate:"omitempty,min=5,max=300"` // 描述
	Parent bson.ObjectId `fill:"iu" json:"parent,omitempty" bson:",omitempty"`                                  // 父节点 ID

	service.Di
}

// 获得 Node 实例
func NewNode(body []byte, typ string) (*Node, error) {
	node := &Node{}
	if err := com.ParseJSON(body, typ, node); err != nil {
		panic(err)
	}

	return node, nil
}

// 添加
func (n *Node) Add(node *Node) error {
	if err := com.NewVali().Struct(node); err != "" {
		return errors.New(err)
	}

	return n.C("node").Insert(node)
}

// 保存
func (n *Node) Save(id bson.ObjectId, node *Node) error {
	if id == "" {
		return errors.New("节点ID不能为空")
	}

	set, err := com.Extract(node, "u")
	if err != nil {
		return err
	}

	return n.C("node").UpdateId(id, M{"$set": set})
}

// 查找所有
func (n *Node) FindAll() (*[]Node, error) {
	node := &[]Node{}
	if err := n.C("node").Find(M{}).All(node); err != nil {
		return nil, err
	}

	return node, nil
}

// 节点ID是否存在
func (n *Node) IdExists(id bson.ObjectId) (bool, error) {
	if c, err := n.C("node").FindId(id).Count(); err != nil {
		return false, err
	} else if c != 1 {
		return false, nil
	}

	return true, nil
}

// 节点名是否存在
func (n *Node) NameExists(name string) (bool, error) {
	if name == "" {
		return false, errors.New("节点名不能为空")
	}

	if c, err := n.C("node").Find(M{"name": name}).Count(); err != nil {
		return false, err
	} else if c < 1 {
		return false, nil
	}

	return true, nil
}

// 是否有子节点存在
func (n *Node) HasChild(id bson.ObjectId) (bool, error) {
	if c, err := n.C("node").Find(M{"parent": id}).Count(); err != nil {
		return false, err
	} else if c < 1 {
		return false, nil
	}

	return true, nil
}

// 节点下是否有主题存在
func (n *Node) HasTopic(id bson.ObjectId) (bool, error) {
	if c, err := n.C("topic").Find(M{"node": id}).Count(); err != nil {
		return false, err
	} else if c < 1 {
		return false, nil
	}

	return true, nil
}

// 通过ID或名称查找
func (n *Node) FindByIdOrName(node *Node) error {
	var q *mgo.Query
	if node.Id != "" {
		q = n.C("node").FindId(node.Id)
	} else if node.Name != "" {
		q = n.C("node").Find(M{"name": node.Name})
	} else {
		return errors.New("ID 或名称不能为空")
	}

	return q.One(node)
}

func (n *Node) RemoveById(id bson.ObjectId) error {
	if id == "" {
		return errors.New("ID 不能为空")
	}

	return n.C("node").RemoveId(id)
}
