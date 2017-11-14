package db

import (
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

	service.Di `json:"-" bson:"-"`
}

// 获得 Bill 实例
func NewBill(body []byte, typ string) (bill *Bill) {
	bill = &Bill{}

	if err := com.ParseJSON(body, typ, bill); err != nil {
		panic(err)
	}

	return
}

// 添加
func (b *Bill) Add(bill *Bill) {
	if err := b.C("bill").Insert(bill); err != nil {
		panic(err.Error())
	}
}

// 通过所属者查找
func (b *Bill) FindByMaster(master bson.ObjectId) (bill []*Bill) {
	if master == "" {
		panic("所属者ID不能为空")
	}

	err := b.C("bill").Find(M{"master": master}).Select(M{"master": 0}).All(&bill)
	if err != nil {
		panic(err.Error())
	}

	return
}
