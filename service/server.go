package service

type Message struct {
	client  *client
	content []byte
}

type Server struct {
	// 已注册的客户
	clients map[*client]bool

	// 来自客户的消息
	broadcast chan *Message

	// 待注册的客户
	register chan *client
	// 待销毁的客户
	unregister chan *client
}

func (s *Server) watch() {
	for {
		select {

		// 注册客户
		case client := <-s.register:
			s.clients[client] = true

		// 销毁客户
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}

		// 处理消息
		case message := <-s.broadcast:
			go func() {
				new(request).handle(message.client, message.content)
			}()
		}
	}
}

func newServer() *Server {
	return &Server{
		clients:    make(map[*client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}
