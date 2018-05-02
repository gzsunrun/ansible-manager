package sockets

// 连接
type hub struct {
	// 链接集合
	connections map[*connection]bool

	// 广播信息队列
	broadcast chan *sendRequest

	// 注册链接队列
	register chan *connection

	// 注销链接队列
	unregister chan *connection
}

type sendRequest struct {
	taskID string
	msg    []byte
}

var h = hub{
	broadcast:   make(chan *sendRequest),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			for c := range h.connections {
				if m.taskID != "" && m.taskID != c.taskID {
					continue
				}

				select {
				case c.send <- m.msg:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

// StartWS start ws listen
func StartWS() {
	go h.run()
}
