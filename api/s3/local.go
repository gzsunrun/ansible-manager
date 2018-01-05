package s3

import (
	"os"
	"os/exec"

	"github.com/gzsunrun/ansible-manager/config"
)

func LocalCopy(object string, save string) error {
	return exec.Command("cp", "-rf", config.Cfg.AnsibleManager.WorkPath+"/repo/"+object, save).Run()
}

func LocalDel(object string) error {
	return os.Remove(config.Cfg.AnsibleManager.WorkPath + "/repo/" + object)
}
