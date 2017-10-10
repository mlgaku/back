package middleware

import (
	"errors"
	. "github.com/mlgaku/back/service"
)

type Auth struct{}

func (*Auth) IsLogin() interface{} {
	return func(req *Request) error {
		return errors.New("balabal")
	}
}
