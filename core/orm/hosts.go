package orm

import (
	"time"

	"github.com/hashwing/log"
)

// Hosts host table
type Hosts struct {
	ID       string    `xorm:"host_id" json:"host_id"`
	UserID   string    `xorm:"user_id" json:"user_id"`
	Alias    string    `xorm:"host_alias" json:"host_alias"`
	HostName string    `xorm:"host_name" json:"host_name"`
	IP       string    `xorm:"host_ip" json:"host_ip"`
	User     string    `xorm:"host_user" json:"host_user"`
	Password string    `xorm:"host_password" json:"host_password"`
	Key      string    `xorm:"host_key" json:"host_key"`
	Status   bool      `xorm:"host_status" json:"host_status"`
	Created  time.Time `xorm:"created" json:"created"`
}

// HostsList host json for output
type HostsList struct {
	ID       string    `xorm:"host_id" json:"host_id"`
	UserID   string    `xorm:"user_id" json:"-"`
	Alias    string    `xorm:"host_alias" json:"host_alias"`
	HostName string    `xorm:"host_name" json:"host_name"`
	IP       string    `xorm:"host_ip" json:"host_ip"`
	User     string    `xorm:"host_user" json:"-"`
	Password string    `xorm:"host_password" json:"-"`
	Key      string    `xorm:"host_key" json:"-"`
	Status   bool      `xorm:"host_status" json:"host_status"`
	Created  time.Time `xorm:"created" json:"created"`
}

// CreateHost insert host into table
func CreateHost(host *Hosts) error {
	if host.Password != "" {
		psw, err := RsaEncrypt([]byte(host.Password))
		if err != nil {
			log.Error(err)
			return err
		}
		host.Password = string(psw)
	}
	if host.Key != "" {
		key, err := RsaEncrypt([]byte(host.Key))
		if err != nil {
			log.Error(err)
			return err
		}
		host.Key = string(key)
	}
	_, err := MysqlDB.Table("ansible_host").Insert(host)
	if err != nil {
		log.Error(err)
	}
	return err
}

// CreateHostList insert hostlist into table
func CreateHostList(host *HostsList) error {
	if host.Password != "" {
		psw, err := RsaEncrypt([]byte(host.Password))
		if err != nil {
			log.Error(err)
			return err
		}
		host.Password = string(psw)
	}
	if host.Key != "" {
		key, err := RsaEncrypt([]byte(host.Key))
		if err != nil {
			log.Error(err)
			return err
		}
		host.Key = string(key)
	}
	_, err := MysqlDB.Table("ansible_host").Insert(host)
	if err != nil {
		log.Error(err)
	}
	return err
}

// UPdateHost update host into table
func UPdateHost(host *Hosts) error {
	if host.Password != "" {
		psw, err := RsaEncrypt([]byte(host.Password))
		if err != nil {
			log.Error(err)
			return err
		}
		host.Password = string(psw)
	}
	if host.Key != "" {
		key, err := RsaEncrypt([]byte(host.Key))
		if err != nil {
			log.Error(err)
			return err
		}
		host.Key = string(key)
	}
	_, err := MysqlDB.Table("ansible_host").Where("host_id=?", host.ID).Update(host)
	if err != nil {
		log.Error(err)
	}
	return err
}
// UPdateNullHost update host enve if filed is null
func UPdateNullHost(host *Hosts) error {
	if host.Password != "" {
		psw, err := RsaEncrypt([]byte(host.Password))
		if err != nil {
			log.Error(err)
			return err
		}
		host.Password = string(psw)
	}
	if host.Key != "" {
		key, err := RsaEncrypt([]byte(host.Key))
		if err != nil {
			log.Error(err)
			return err
		}
		host.Key = string(key)
	}
	_, err := MysqlDB.Table("ansible_host").Where("host_id=?", host.ID).
		Cols("host_alias", "host_name", "host_ip", "host_user", "host_password", "host_key", "host_status").
		Update(host)
	if err != nil {
		log.Error(err)
	}
	return err
}

// UPdateAuthHost update and auth host
func UPdateAuthHost(host *Hosts) error {
	var err error
	if host.Password != "" {
		psw, err := RsaEncrypt([]byte(host.Password))
		if err != nil {
			log.Error(err)
			return err
		}
		host.Password = string(psw)
	}
	if host.Key != "" {
		key, err := RsaEncrypt([]byte(host.Key))
		if err != nil {
			log.Error(err)
			return err
		}
		host.Key = string(key)
	}
	if host.HostName == "" {
		_, err = MysqlDB.Table("ansible_host").Where("host_id=?", host.ID).
			Cols("host_user", "host_password", "host_key").
			Update(host)
	} else {
		_, err = MysqlDB.Table("ansible_host").Where("host_id=?", host.ID).
			Cols("host_name", "host_user", "host_password", "host_key").
			Update(host)
	}
	if err != nil {
		log.Error(err)
	}
	return err
}

// FindHosts find hosts
func FindHosts(uid string, hosts *[]HostsList) error {
	err := MysqlDB.Table("ansible_host").Where("user_id=?", uid).Find(hosts)
	if err != nil {
		log.Error(err)
		return err
	}
	for i, h := range *hosts {
		psw, err := RsaDecrypt(h.Password)
		if err == nil {
			(*hosts)[i].Password = psw
		}
		key, err := RsaDecrypt(h.Key)
		if err == nil {
			(*hosts)[i].Key = key
		}
	}
	return nil
}

// GetHost get host
func GetHost(hostID string, host interface{}) (bool, error) {
	res, err := MysqlDB.Table("ansible_host").Where("host_id=?", hostID).Get(host)
	if err != nil {
		log.Error(err)
	}

	if ehost, ok := host.(*Hosts); ok {
		psw, err := RsaDecrypt(ehost.Password)
		if err == nil {
			ehost.Password = psw
		}
		key, err := RsaDecrypt(ehost.Key)
		if err == nil {

			ehost.Key = key
		}
		host = ehost
	}
	if ehost, ok := host.(*HostsList); ok {
		psw, err := RsaDecrypt(ehost.Password)
		if err == nil {
			ehost.Password = psw
		}
		key, err := RsaDecrypt(ehost.Key)
		if err == nil {

			ehost.Key = key
		}
		host = ehost
	}
	return res, err
}

// DelHost delete host
func DelHost(hid string) error {
	host := new(Hosts)
	_, err := MysqlDB.Table("ansible_host").Where("host_id=?", hid).Delete(host)
	if err != nil {
		log.Error(err)
	}
	return err
}
