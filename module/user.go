package module

import (
	"encoding/json"
	. "github.com/mlgaku/back/types"
	//"gopkg.in/mgo.v2/bson"
	"github.com/mlgaku/back/common"
)

type User struct {
	Name     string `json:"name" validate:"required,min=4,max=15,alphanum"`
	Email    string `json:"email" validate:"required,min=8,max=30,email"`
	Password string `json:"password" validate:"required,min=8,max=20,alphanum"`
}

func (*User) Reg(db *Database, req *Request) Value {

	user := &User{}
	json.Unmarshal(req.Body, user)

	if err := common.NewValidator().Struct(user); err != "" {
		return err
	}
	return "验证通过"

	//c := db.Session.DB("test").C("people")
	//
	//err := c.Insert(&Person{"Ale", "+55 53 8116 9639"}, &Person{"Cla", "+55 53 8402 8510"})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//result := Person{}
	//err = c.Find(bson.M{"name": "Ale"}).One(&result)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("Phone:", result.Phone)
	return ""
}
