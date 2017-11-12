package db

import (
	"errors"
	com "github.com/mlgaku/back/common"
	"github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Bill struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"` // ID
	Msg    string        `json:"msg" bson:",omitempty"`   // 提示信息
	Type   uint64        `json:"type" bson:",minsize"`    // 类型(1.新主题 2.回复主题 3.签到)
	Date   time.Time     `json:"date"`                    // 日期
	Number int64         `json:"number"`                  // 数量
	Master bson.ObjectId `json:"master,omitempty"`        // 所属者ID

	service.Di
}

// 获得 Bill 实例
func NewBill(body []byte, typ string) *Bill {
	bill := &Bill{}
	if err := com.ParseJSON(body, typ, bill); err != nil {
		panic(err)
	}

	return bill
}

// 添加
func (b *Bill) Add(bill *Bill) error {
	return b.C("bill").Insert(bill)
}

// 通过所属者查找
func (b *Bill) FindByMaster(master bson.ObjectId) (*[]Bill, error) {
	if master == "" {
		return nil, errors.New("所属者ID不能为空")
	}

	bill := &[]Bill{}
	err := b.C("bill").Find(M{"master": master}).Select(M{"master": 0}).All(bill)
	if err != nil {
		return nil, err
	}

	return bill, nil
}
