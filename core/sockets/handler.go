package sockets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gzsunrun/ansible-manager/core/output"
	"github.com/hashwing/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	// 写信息时等待时间
	writeWait = 10 * time.Second

	// 心跳超时时间
	pongWait = 60 * time.Second

	// 发送心跳包，必须比心跳超时时间短
	pingPeriod = (pongWait * 9) / 10

	// 信息缓存大小
	maxMessageSize = 512
)

type connection struct {
	ws     *websocket.Conn
	send   chan []byte
	taskID string
}

// readPump 从连接中读信息
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		fmt.Print(string(message))

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Error("error:", err)
			}
			break
		}
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump 从连接中写信息
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

//Handler 建立websocket 连接
func Handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	taskID := r.FormValue("task_id")
	if taskID == "" {
		log.Error(err)
		return
	}
	c := &connection{
		send:   make(chan []byte, 256),
		ws:     ws,
		taskID: taskID,
	}

	h.register <- c

	go c.writePump()
	lo, err := output.NewLogOutput(taskID)
	if err != nil {
		log.Error(err)
	}
	data, err := lo.Read()
	if err != nil {
		log.Error(err)
	}
	for _, msg := range data {
		b, err := json.Marshal(&map[string]interface{}{
			"type":    "log",
			"output":  msg,
			"task_id": taskID,
		})
		if err != nil {
			log.Error(err)
			continue
		}
		err = c.ws.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			log.Error(err)
		}
	}

	go Client(taskID)
	c.readPump()
}

// Message 往广播信息队列中写信息
func Message(taskID string, message []byte) {
	h.broadcast <- &sendRequest{
		taskID: taskID,
		msg:    message,
	}
}

// CloseConn 断开连接
func CloseConn(taskID string) {
	for c := range h.connections {
		if c.taskID == taskID {
			h.unregister <- c
			c.ws.Close()
		}
	}
}
