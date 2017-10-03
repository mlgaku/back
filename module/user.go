package module

import (
	"encoding/json"
	"fmt"
	. "github.com/mlgaku/back/types"
	//"gopkg.in/go-playground/validator.v9"
	//"gopkg.in/mgo.v2/bson"
	//"log"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (*User) Reg(db *Database, req *Request) Value {

	user := &User{}
	json.Unmarshal(req.Body, &user)
	fmt.Println(user)

	//validate := validator.New()
	//
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
