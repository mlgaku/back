package module

import (
	"fmt"
	"github.com/mlgaku/back/common"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Person struct {
	Name  string
	Phone string
}

type User struct{}

func (*User) Reg(db *common.Database) common.Value {
	c := db.Session.DB("test").C("people")

	err := c.Insert(&Person{"Ale", "+55 53 8116 9639"}, &Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
	return ""
}

func (*User) Test() common.Value {
	return map[string]string{"name": "yazi", "age": "17"}

	//return &struct {
	//	Name  string `json:"name"`
	//	Phone string `json:"phone"`
	//}{
	//	Name:  "yazi",
	//	Phone: "110",
	//}
}
