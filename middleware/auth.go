package middleware

import (
	. "github.com/mlgaku/back/service"
)

type Auth struct{}

func (*Auth) IsLogin() interface{} {
	return func(req *Request) error {
		return nil
	}
}
