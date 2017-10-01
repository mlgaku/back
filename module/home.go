package module

import (
	"fmt"
	. "github.com/sxyazi/maile/common"
	"math/rand"
)

type Home struct{}

// 测试
func (*Home) Hello(req *Request, res *Response) Value {
	for i := 0; i < 10; i++ {
		res.Write("ddd")
	}

	fmt.Println(req.BodyStr)

	return fmt.Sprintf("[%g] hello world", rand.Float64())
}
