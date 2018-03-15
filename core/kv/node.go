package kv

import (
	"os"
	"net"
	"errors"
	"strconv"
	"encoding/json"

	"github.com/satori/go.uuid"
)

type Node struct{
	NodeID	 	string	`json:"node_id"`
	IP 	 		string	`json:"node_ip"`
	TTL  		int64	`json:"node_ttl"`
	PID  		string	`json:"node_pid"`
	Port 		int		`json:"node_port"`
	Path 		string	`json:"node_path"`
	Worker 		bool	`json:"node_worker"`
	Master		bool	`json:"node_master"`
}

func NewNode(ttl int64,port int,path string,worker,master bool)(*Node,error){
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return &Node{
					NodeID:uuid.Must(uuid.NewV4()).String(),
					IP:ipnet.IP.String(),
					TTL:ttl,
					Port:port,
					Worker:worker,
					Master:master,
					Path:path,
					PID:strconv.Itoa(os.Getpid()),
				},nil
			}
		}
	}
	return nil, errors.New("not found ip")
}

func (n *Node)ID()string{
	return n.NodeID
}

func (n *Node)OutTime()int64{
	return n.TTL
}

func (n *Node)String()string{
	data,_:=json.Marshal(n)
	return string(data)
}