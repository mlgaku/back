package types

type Pubsub struct {
	Publish   func(id string)
	AddSub    func(id string, res *Response)
	RemoveSub func(id string, res *Response)
}
