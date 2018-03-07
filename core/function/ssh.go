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
		res:=AuthPassword(host)
		if res!="success"{
			if host.Key!=""{
				return AuthKey(host)
			}
		}
		return res
	}
	if host.Key!=""{
		return AuthKey(host)
	}
	return "fail"
} 

func AuthPassword(host orm.HostsList)string{
	authMethods := []ssh.AuthMethod{}

    keyboardInteractiveChallenge := func(
        user,
        instruction string,
        questions []string,
        echos []bool,
    ) (answers []string, err error) {
        if len(questions) == 0 {
            return []string{}, nil
        }
        return []string{host.Password}, nil
    }
	
	authMethods = append(authMethods, ssh.KeyboardInteractive(keyboardInteractiveChallenge))
    authMethods = append(authMethods, ssh.Password(host.Password))

	c, err := ssh.Dial("tcp", host.IP+":22", &ssh.ClientConfig{
		User: host.User,
		Auth: authMethods,
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

func AuthKey(host orm.HostsList)string{
	signer,err:=ssh.ParsePrivateKey([]byte(host.Key))
		if err!=nil{
			log.Error(err)
			return "fail"
		}
		if host.User==""{
			host.User="root"
		}
		c, err := ssh.Dial("tcp", host.IP+":22", &ssh.ClientConfig{
			User: host.User,
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

func AuthKeyByHost(host orm.Hosts)string{
	signer,err:=ssh.ParsePrivateKey([]byte(host.Key))
		if err!=nil{
			log.Error(err)
			return "fail"
		}
		c, err := ssh.Dial("tcp", host.IP+":22", &ssh.ClientConfig{
			User: host.User,
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
