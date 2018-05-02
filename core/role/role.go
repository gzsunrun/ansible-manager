package role

import (
	"github.com/gzsunrun/ansible-manager/core/config"
	"github.com/gzsunrun/ansible-manager/core/kv"
)

var node = func(n kv.Node, p bool) {}
var task = func(t kv.Task, p bool) {}
var sche = func(t kv.Task, p bool) {}

// Run run role
func Run() {
	if config.Cfg.Common.Master {
		MasterSet()
	}
	if config.Cfg.Common.Worker {
		WorkerSet()
	}
	kv.DefaultClient.SetCall(node, task, sche)
}
