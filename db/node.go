package db

import (
	"errors"
	com "github.com/mlgaku/back/common"
	. "github.com/mlgaku/back/service"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Node struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name   string        `json:"name" validate:"required,max=30,alphanum"`
	Title  string        `json:"title" validate:"required,max=30"`
	Parent bson.ObjectId `json:"parent,omitempty"`
}

// 添加
func (n *Node) Add(db *Database, node *Node) error {
	if err := com.NewVali().Struct(node); err != "" {
		return errors.New(err)
	}

	// 检查父节点
	if node.Parent != "" {
		if b, err := n.IdExists(db, node.Parent); err != nil {
			return err
		} else if !b {
			return errors.New("父节点不存在")
		}
	}

	if err := db.C("node").Insert(node); err != nil {
		return err
	}

	return nil
}

// 查询所有
func (*Node) FindAll(db *Database) (*[]Node, error) {
	node := &[]Node{}
	if err := db.C("node").Find(bson.M{}).All(node); err != nil {
		return nil, err
	}

	return node, nil
}

// 节点ID是否存在
func (*Node) IdExists(db *Database, id bson.ObjectId) (bool, error) {
	if c, err := db.C("node").FindId(id).Count(); err != nil {
		return false, err
	} else if c != 1 {
		return false, nil
	}

	return true, nil
}

// 节点名是否存在
func (*Node) NameExists(db *Database, name string) (bool, error) {
	if name == "" {
		return false, errors.New("节点名不能为空")
	}

	if c, err := db.C("node").Find(bson.M{"name": name}).Count(); err != nil {
		return false, err
	} else if c < 1 {
		return false, nil
	}

	return true, nil
}

// 是否有子节点存在
func (*Node) HasChild(db *Database, id bson.ObjectId) (bool, error) {
	if c, err := db.C("node").Find(bson.M{"parent": id}).Count(); err != nil {
		return false, err
	} else if c < 1 {
		return false, nil
	}

	return true, nil
}

// 通过ID或名称查询
func (*Node) FindByIdOrName(db *Database, node *Node) error {
	var q *mgo.Query
	if node.Id != "" {
		q = db.C("node").FindId(node.Id)
	} else if node.Name != "" {
		q = db.C("node").Find(bson.M{"name": node.Name})
	} else {
		return errors.New("ID 或名称不能为空")
	}

	if err := q.One(node); err != nil {
		return err
	}
	return nil
}

func (n *Node) RemoveById(db *Database, id bson.ObjectId) error {
	if id == "" {
		return errors.New("ID 不能为空")
	}

	if b, err := n.HasChild(db, id); err != nil {
		return err
	} else if b {
		return errors.New("删除失败: 该节点下有子节点存在")
	}

	if err := db.C("node").RemoveId(id); err != nil {
		return err
	}

	return nil
}
