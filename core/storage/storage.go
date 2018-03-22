package storage

import (
	"github.com/gzsunrun/ansible-manager/core/config"
)

type StorageParse struct {
	RemotePath string
	LocalPath  string
}

type FileStorage interface {
	Put(*StorageParse) error
	Get(*StorageParse) error
	Delete(*StorageParse) error
	Share(*StorageParse) (string, error)
}

var Storage FileStorage

func SetStorage() error {
	var err error
	if config.Cfg.S3.Enable {
		Storage, err = NewS3Storage(config.Cfg.S3.S3URL,
			config.Cfg.S3.S3Key,
			config.Cfg.S3.S3Secret,
			config.Cfg.S3.BucketName,
		)
	} else if config.Cfg.Git.Enable {
		Storage = NewGit()
	} else {
		Storage, err = NewLocalStorage(config.Cfg.LocalStorage.Path)
	}
	return err
}

func NewStorage() (FileStorage, error) {
	if config.Cfg.S3.Enable {
		return NewS3Storage(config.Cfg.S3.S3URL,
			config.Cfg.S3.S3Key,
			config.Cfg.S3.S3Secret,
			config.Cfg.S3.BucketName,
		)
	}

	return NewLocalStorage(config.Cfg.LocalStorage.Path)
}
