package iaas

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"git.gzsunrun.cn/sunruniaas/sunruniaas-core/common"
	"git.gzsunrun.cn/sunruniaas/sunruniaas-core/rpc"
	"github.com/gzsunrun/ansible-manager/core/orm"
	log "github.com/astaxie/beego/logs"
)

type ProjectHosts struct{
	Uuid string
}
func GetProjectHosts(uuid string)(*[]orm.HostsList,error){
	hosts:=make([]orm.HostsList,0)
	client1 := rpc.NewMbusClientWithConsul()
	defer client1.Close()
	reply := new(rpc.ProjectMetaReply)
	err := client1.Call(context.Background(),rpc.C_MC_RPC_PROJ_GET_META, uuid, reply)
	if err!=nil{
		log.Error(err)
		return nil,err
	}
	v,err:=json.Marshal(reply.Meta.Resources)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var data []ProjectHosts
	err=json.Unmarshal(v,&data)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	client2 := rpc.NewMbusClientWithConsul()
	defer client2.Close()
	for _, h := range data {
		reply := new(rpc.InstanceMetaReply)
		err := client2.Call(context.Background(),
			rpc.C_MC_RPC_VM_GET_METADATA, h.Uuid,reply)

		if err != nil {
			return nil, err
		}
		if len(reply.Meta.VnetIfaces) == 0 {
			return nil, errors.New("hosts not found netIfaces")
		}
		host:=orm.HostsList{}
		host.Alias=reply.Meta.Name
		host.ID=reply.Meta.Uuid
		for _, v := range reply.Meta.Tags {
			if strings.HasPrefix(v.Tag, common.C_TAG_META_HOSTNAME_PREFIX) {
				host.Name = strings.Replace(v.Tag, common.C_TAG_META_HOSTNAME_PREFIX, "", -1)
				continue
			}
			if strings.HasPrefix(v.Tag, common.C_TAG_META_KEYPAIR_PREFIX) {
				keyStr := strings.Replace(v.Tag, common.C_TAG_META_KEYPAIR_PREFIX, "", -1)
				keys := strings.Split(keyStr, ":")
				key, err := base64.StdEncoding.DecodeString(keys[0])
				if err != nil {
					return nil, err
				}
				host.Key = string(key)
			}
		}
		for _, v := range reply.Meta.VnetIfaces {
			if v.IsDefaultNic {
				host.IP = v.IpStrval
			}
		}
		hosts=append(hosts,host)
	}
	return &hosts,nil
}
