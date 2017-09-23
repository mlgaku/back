package module

import (
	"fmt"
	. "github.com/sxyazi/maile/common"
	"math/rand"
)

type Home struct{}

// 测试
func (*Home) Hello() Value {
	return fmt.Sprintf("[%g] hello world", rand.Float64())
}
