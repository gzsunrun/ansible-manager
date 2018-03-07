package storage

import(
	"github.com/gzsunrun/ansible-manager/core/config"
)

type StorageParse struct{
	RemotePath string
	LocalPath string
}

type FileStorage interface{
	Put(*StorageParse)error
	Get(*StorageParse)error
	Delete(*StorageParse)error
	Share(*StorageParse)(string,error)
}

var Storage FileStorage

func SetStorage()error{
	var err error
	if config.Cfg.Ansible.S3Status{
		Storage,err = NewS3Storage(config.Cfg.Ansible.S3URL,
			config.Cfg.Ansible.S3Key,
			config.Cfg.Ansible.S3Secret,
			config.Cfg.Ansible.BucketName,
		)
	}else{
		Storage,err = NewLocalStorage(config.Cfg.Ansible.WorkPath+"/repo/")
	}
	return err
}

func NewStorage()(FileStorage,error){
	if config.Cfg.Ansible.S3Status{
		return NewS3Storage(config.Cfg.Ansible.S3URL,
			config.Cfg.Ansible.S3Key,
			config.Cfg.Ansible.S3Secret,
			config.Cfg.Ansible.BucketName,
		)
	}

	return NewLocalStorage(config.Cfg.Ansible.WorkPath+"/repo/")
}