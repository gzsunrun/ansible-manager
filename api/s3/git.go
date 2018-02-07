package s3

import (
	"os/exec"

	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/config"
)

func GitClone(url,path string)error{
	cmd:=exec.Command("git","clone",url,path+"_dir")
	cmd.Dir=config.Cfg.AnsibleManager.WorkPath
	err:=cmd.Run()
	if err!=nil{
		logs.Error("clone:",err)
		return err
	}
	defer exec.Command("rm","-rf",config.Cfg.AnsibleManager.WorkPath+"/"+path+"_dir").Run()

	cmd =exec.Command("tar","zcvf",config.Cfg.AnsibleManager.WorkPath+"/repo/"+path,"./")
	cmd.Dir=config.Cfg.AnsibleManager.WorkPath+"/"+path+"_dir"
	return cmd.Run()
}