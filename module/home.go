package module

import (
	"fmt"
	. "github.com/mlgaku/back/common"
	"math/rand"
)

type Home struct{}

// 测试
func (*Home) Hello(req *Request, res *Response) Value {
	return &struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "yazi",
		Age:  17,
	}

	for i := 0; i < 10; i++ {
		res.Write("ddd")
	}

	fmt.Println(req.BodyStr)

	return fmt.Sprintf("[%g] hello world", rand.Float64())
}
