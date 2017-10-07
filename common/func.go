package common

import (
	"strconv"
	"strings"
)

func IPAddr(a string) (ip string, port int) {
	s := strings.SplitN(a, ":", 2)
	if len(s) < 2 {
		s = append(s, "0")
	}

	ip = s[0]
	port, _ = strconv.Atoi(s[1])
	return
}
