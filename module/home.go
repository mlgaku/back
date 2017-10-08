package module

import (
	"fmt"
	. "github.com/mlgaku/back/service"
	. "github.com/mlgaku/back/types"
	"math/rand"
)

type Home struct{}

// 测试
func (*Home) Hello(req *Request, res *Response) Value {
	for i := 0; i < 10; i++ {
		res.Write([]byte("ddd"))
	}

	return fmt.Sprintf("[%g] hello world", rand.Float64())
}
