package function

import (
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/paulstuart/ping"
)

// SshDail ssh auth
func SshDail(host orm.HostsList) string {
	netStatus := ping.Ping(host.IP, 5)
	if !netStatus {
		return "fail"
	}
	if host.Password != "" {
		res := AuthPassword(host)
		if res != "success" {
			if host.Key != "" {
				return AuthKey(host)
			}
		}
		return res
	}
	if host.Key != "" {
		return AuthKey(host)
	}
	return "fail"
}

// AuthPassword ssh auth by password
func AuthPassword(host orm.HostsList) string {
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

	c, err := ssh.Dial("tcp", host.IP+":"+host.Port, &ssh.ClientConfig{
		User:    host.User,
		Auth:    authMethods,
		Timeout: 3 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})

	if err != nil {
		if strings.Contains(err.Error(), "unable to authenticate") {
			return "auth"
		}
		return "fail"
	}
	defer c.Close()
	return "success"
}

// AuthKey ssh auth by key
func AuthKey(host orm.HostsList) string {
	signer, err := ssh.ParsePrivateKey([]byte(host.Key))
	if err != nil {
		return "fail"
	}
	if host.User == "" {
		host.User = "root"
	}
	c, err := ssh.Dial("tcp", host.IP+":"+host.Port, &ssh.ClientConfig{
		User:    host.User,
		Auth:    []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Timeout: 3 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})

	if err != nil {
		if strings.Contains(err.Error(), "unable to authenticate") {
			return "auth"
		}
		return "fail"
	}
	defer c.Close()
	return "success"
}

// AuthKeyByHost ssh auth by key with Hosts
func AuthKeyByHost(host orm.Hosts) string {
	signer, err := ssh.ParsePrivateKey([]byte(host.Key))
	if err != nil {
		return "fail"
	}
	c, err := ssh.Dial("tcp", host.IP+":"+host.Port, &ssh.ClientConfig{
		User:    host.User,
		Auth:    []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Timeout: 3 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})

	if err != nil {
		fmt.Println(host.IP, host.Password, err)
		if strings.Contains(err.Error(), "unable to authenticate") {
			return "auth"
		}
		return "fail"
	}
	defer c.Close()
	return "success"
}
