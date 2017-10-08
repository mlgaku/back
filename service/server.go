package service

type (
	Server struct {
		// 已注册的客户
		clients map[*Client]bool

		// 来自客户的消息
		broadcast chan *message

		// 待注册的客户
		register chan *Client
		// 待销毁的客户
		unregister chan *Client
	}

	message struct {
		client  *Client
		content []byte
	}
)

// 监听
func (s *Server) Watch(han func(*Client, []byte)) {
	for {
		select {

		// 注册客户
		case client := <-s.register:
			s.clients[client] = true

		// 销毁客户
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.Send)
			}

		// 处理消息
		case message := <-s.broadcast:
			go func() {
				han(message.client, message.content)
			}()

		}
	}
}

// 获得 Server 实例
func NewServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}
