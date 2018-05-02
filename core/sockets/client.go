package sockets

import (
	"fmt"
	"time"

	log "github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/gzsunrun/ansible-manager/core/kv"
)

// WsChan ws client chan
var WsChan = make(map[string]chan bool)

// Client create a client by task id
func Client(taskID string) {
	if WsChan[taskID] != nil {
		return
	}
	WsChan[taskID] = make(chan bool)
	defer func() {
		WsChan[taskID] = nil
	}()
	nodeID := kv.DefaultClient.GetStorage().Tasks[taskID].NodeID
	if nodeID == kv.DefaultClient.LocalNode().ID() {
		return
	}
	ip := kv.DefaultClient.GetStorage().Nodes[nodeID].IP
	port := kv.DefaultClient.GetStorage().Nodes[nodeID].Port
	path := kv.DefaultClient.GetStorage().Nodes[nodeID].Path
	if ip == "" {
		return
	}
	addr := fmt.Sprintf("ws://%s:%d%s%s", ip, port, path, "?task_id="+taskID)
	log.Info(addr)
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		log.Error("dial:", err)
		return
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Error("read:", err)
				return
			}
			Message(taskID, message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-WsChan[taskID]:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Error("write close:", err)
				return
			}
			// close  local connect
			CloseConn(taskID)
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// StopClient close client
func StopClient(taskID string) {
	if WsChan[taskID] == nil {
		return
	}
	WsChan[taskID] <- true
}
