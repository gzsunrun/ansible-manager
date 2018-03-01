package function

import (
	"time"
	"golang.org/x/crypto/ssh"
	"strings"

	log "github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/core/orm"
	
)


func SshDail(host orm.HostsList)string{
	if host.Password!=""{
		c, err := ssh.Dial("tcp", host.IP+":22", &ssh.ClientConfig{
			User: host.User,
			Auth: []ssh.AuthMethod{ssh.Password("sunrunvas")},
			Timeout:3*time.Second,
		})
		
		if err!=nil{
			if strings.Contains(err.Error(),"unable to authenticate"){
				return "auth"
			}
			return "fail"
		}
		defer c.Close()
		return "success"
	}
	if host.Key!=""{
		signer,err:=ssh.ParsePrivateKey([]byte(host.Key))
		if err!=nil{
			log.Error(err)
			return "fail"
		}
		c, err := ssh.Dial("tcp", host.IP+":22", &ssh.ClientConfig{
			User: "root",
			Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
			Timeout:3*time.Second,
		})
		
		if err!=nil{
			if strings.Contains(err.Error(),"unable to authenticate"){
				return "auth"
			}
			return "fail"
		}
		defer c.Close()
		return "success"
	}
	return "fail"
} 