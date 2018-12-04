package kv

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashwing/log"
	"go.etcd.io/etcd/clientv3"
)

var (
	rootPath       string
	tasksPath      string
	nodesPath      string
	connectsPath   string
	locksPath      string
	timerPath      string
	schedulersPath string
)

// EtcdClient etcd client
type EtcdClient struct {
	client     *clientv3.Client
	reqTimeout time.Duration
	lID        clientv3.LeaseID
	Storage    *LocalMenory
	Node       *Node
	nodeCall   func(Node, bool)
	taskCall   func(Task, bool)
	scheCall   func(Task, bool)
}

// NewEtcd new etcd client
func NewEtcd(endpoints []string, ttl int64, port int, path string, worker, master bool) (*EtcdClient, error) {
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(ttl) * time.Second,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	node, err := NewNode(ttl, port, path, worker, master)
	if err != nil {
		return nil, err
	}
	storage := &LocalMenory{
		Tasks: make(map[string]Task),
		Nodes: make(map[string]Node),
		Lock:  new(sync.Mutex),
	}
	rootPath = "ansible/"
	schedulersPath = "ansible/schedulers/"
	tasksPath = "ansible/tasks/"
	nodesPath = "ansible/nodes/"
	connectsPath = "ansible/connects/"
	locksPath = "ansible/locks/"
	return &EtcdClient{
		client:     client,
		reqTimeout: time.Duration(ttl) * time.Second,
		Node:       node,
		Storage:    storage,
	}, nil
}

// LocalNode get this node info
func (ec *EtcdClient) LocalNode() *Node {
	return ec.Node
}

// DelAllTask del all task
func (ec *EtcdClient) DelAllTask() error {
	_, err := ec.Delete(tasksPath, clientv3.WithPrefix())
	return err
}

// AddTask add task
func (ec *EtcdClient) AddTask(task Task) error {
	key := fmt.Sprint(tasksPath, task.ID)
	val, err := json.Marshal(task)
	if err != nil {
		return err
	}

	lgr, err := ec.Grant(5)
	if err != nil {
		return err
	}
	res, err := ec.GetLock(task.ID, lgr.ID)
	if err != nil {
		return err
	}
	if res {
		_, err = ec.Put(key, string(val))
		ec.DeleteScheduler(task.ID)
	}

	return err
}

// DeleteTask delete task
func (ec *EtcdClient) DeleteTask(tid string) error {
	_, err := ec.Delete(tasksPath + tid)
	return err
}

// GetStorage get storage
func (ec *EtcdClient) GetStorage() *LocalMenory {
	return ec.Storage
}

// AddScheduler add  scheduler
func (ec *EtcdClient) AddScheduler(task Task) error {
	key := fmt.Sprint(schedulersPath, task.ID)
	val, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, err = ec.Put(key, string(val))
	return err
}

// DeleteScheduler delete scheduler
func (ec *EtcdClient) DeleteScheduler(tid string) error {
	_, err := ec.Delete(schedulersPath + tid)
	return err
}

// SetCall set call function
func (ec *EtcdClient) SetCall(node func(Node, bool), task, sche func(Task, bool)) {
	ec.taskCall = task
	ec.nodeCall = node
	ec.scheCall = sche
	go ec.WatchData()
}

// SyncData sync data from etcd
func (ec *EtcdClient) SyncData() error {
	var nodes []Node
	err := ec.FindNodes(&nodes)
	if err != nil {
		return err
	}
	for _, v := range nodes {
		ec.Storage.Nodes[v.ID()] = v
	}

	var tasks []Task
	err = ec.FindTasks(&tasks)
	if err != nil {
		return err
	}

	for _, v := range tasks {
		ec.Storage.Tasks[v.ID] = v
	}

	return nil

}

// FindTasks find all tasks
func (ec *EtcdClient) FindTasks(tasks interface{}) error {
	grsp, err := ec.Get(tasksPath, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	str := "["
	for i, v := range grsp.Kvs {
		str += string(v.Value)
		if i < len(grsp.Kvs)-1 {
			str += ","
		}
	}
	str += "]"
	err = json.Unmarshal([]byte(str), tasks)
	return err
}

// RegNode reg node into etcd
func (ec *EtcdClient) RegNode() error {
	lrsp, err := ec.Grant(ec.Node.OutTime() + 2)
	ec.lID = lrsp.ID
	_, err = ec.Put(nodesPath+ec.Node.ID(), ec.Node.String(), clientv3.WithLease(lrsp.ID))
	return err
}

// DelNode delete node
func (ec *EtcdClient) DelNode() error {
	_, err := ec.Delete(nodesPath + ec.Node.ID())
	return err
}

// FindNodes find all nodes
func (ec *EtcdClient) FindNodes(nodes interface{}) error {
	grsp, err := ec.Get(nodesPath, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	str := "["
	for i, v := range grsp.Kvs {
		str += string(v.Value)
		if i < len(grsp.Kvs)-1 {
			str += ","
		}
	}
	str += "]"
	err = json.Unmarshal([]byte(str), nodes)
	return err
}

// GetNode get node
func (ec *EtcdClient) GetNode(node interface{}, id string) error {
	grsp, err := ec.Get(nodesPath + id)
	if err != nil {
		return err
	}

	for _, v := range grsp.Kvs {
		err = json.Unmarshal(v.Value, node)
		return err
	}
	return nil
}

// WatchData watch key when put or delete
func (ec *EtcdClient) WatchData() {
	wch := ec.Watch(rootPath, clientv3.WithPrefix())
	for {
		select {
		case c := <-wch:
			for _, e := range c.Events {
				ec.Storage.Lock.Lock()

				//nodes put/del
				if strings.HasPrefix(string(e.Kv.Key), nodesPath) {
					keys := strings.Split(string(e.Kv.Key), "/")
					if e.Type == clientv3.EventTypeDelete {
						go ec.nodeCall(ec.Storage.Nodes[keys[2]], false)
						delete(ec.Storage.Nodes, keys[2])
					}
					if e.Type == clientv3.EventTypePut {
						var node Node
						err := json.Unmarshal(e.Kv.Value, &node)
						if err != nil {
							log.Error(err)
						}
						ec.Storage.Nodes[keys[2]] = node
						ec.nodeCall(node, true)

					}
				}

				//tasks put/del
				if strings.HasPrefix(string(e.Kv.Key), tasksPath) {
					keys := strings.Split(string(e.Kv.Key), "/")
					if e.Type == clientv3.EventTypeDelete {
						go ec.taskCall(ec.Storage.Tasks[keys[2]], false)
						delete(ec.Storage.Tasks, keys[2])
					}
					if e.Type == clientv3.EventTypePut {
						var task Task
						err := json.Unmarshal(e.Kv.Value, &task)
						if err != nil {
							log.Error(err)
						}
						ec.Storage.Tasks[keys[2]] = task
						ec.taskCall(task, true)
					}
				}

				//scheduler put/del
				if strings.HasPrefix(string(e.Kv.Key), schedulersPath) {

					if e.Type == clientv3.EventTypePut {
						var task Task
						err := json.Unmarshal(e.Kv.Value, &task)
						if err != nil {
							log.Error(err)
						}
						log.Info(task)
						ec.scheCall(task, true)
					}
				}
				ec.Storage.Lock.Unlock()
			}
		}
	}
}

// KeepNode keep node
func (ec *EtcdClient) KeepNode() {
	duration := time.Duration(ec.Node.TTL) * time.Second
	timer := time.NewTimer(duration)
	for {
		select {
		case <-timer.C:
			if ec.lID > 0 {
				_, err := ec.KeepAliveOnce(ec.lID)
				if err == nil {
					timer.Reset(duration)
					continue
				}
				ec.lID = 0
			}

			if err := ec.RegNode(); err != nil {
				log.Error("%s connect err: %s, try to reset after %d seconds...", ec.Node.ID(), err.Error(), ec.Node.OutTime())
			} else {
				log.Info("%s connect success,lid:%x", ec.Node.ID(), ec.lID)
			}
			timer.Reset(duration)
		}
	}
}

// Put etcd put
func (ec *EtcdClient) Put(key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext(ec)
	defer cancel()
	return ec.client.Put(ctx, key, val, opts...)
}

// Get etcd get
func (ec *EtcdClient) Get(key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext(ec)
	defer cancel()
	return ec.client.Get(ctx, key, opts...)
}

// Delete etcd delete
func (ec *EtcdClient) Delete(key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext(ec)
	defer cancel()
	return ec.client.Delete(ctx, key, opts...)
}

// Grant etcd grant
func (ec *EtcdClient) Grant(ttl int64) (*clientv3.LeaseGrantResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext(ec)
	defer cancel()
	return ec.client.Grant(ctx, ttl)
}

// Watch etcd watch
func (ec *EtcdClient) Watch(key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return ec.client.Watch(context.Background(), key, opts...)
}

// KeepAliveOnce keepalive one time
func (ec *EtcdClient) KeepAliveOnce(id clientv3.LeaseID) (*clientv3.LeaseKeepAliveResponse, error) {
	ctx, cancel := NewEtcdTimeoutContext(ec)
	defer cancel()
	return ec.client.KeepAliveOnce(ctx, id)
}

// GetLock get lock
func (ec *EtcdClient) GetLock(key string, id clientv3.LeaseID) (bool, error) {
	key = locksPath + key
	ctx, cancel := NewEtcdTimeoutContext(ec)
	resp, err := ec.client.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "", clientv3.WithLease(id))).
		Commit()
	cancel()

	if err != nil {
		return false, err
	}

	return resp.Succeeded, nil
}

// DelLock delete lock
func (ec *EtcdClient) DelLock(key string) error {
	_, err := ec.Delete(locksPath + key)
	return err
}

// etcdTimeoutContext etcd timeout context
type etcdTimeoutContext struct {
	context.Context

	etcdEndpoints []string
}

// Err err
func (c *etcdTimeoutContext) Err() error {
	err := c.Context.Err()
	if err == context.DeadlineExceeded {
		err = fmt.Errorf("%s: etcd(%v) lost",
			err, c.etcdEndpoints)
	}
	return err
}
// NewEtcdTimeoutContext new etcd timeout context
func NewEtcdTimeoutContext(c *EtcdClient) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), c.reqTimeout)
	etcdCtx := &etcdTimeoutContext{}
	etcdCtx.Context = ctx
	etcdCtx.etcdEndpoints = c.client.Endpoints()
	return etcdCtx, cancel
}
