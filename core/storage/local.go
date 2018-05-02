package storage

import (
	log "github.com/astaxie/beego/logs"
	"io"
	"os"
)

// LocalStorage local repo
type LocalStorage struct {
	LocalPath string
}

// NewLocalStorage new local storage
func NewLocalStorage(path string) (*LocalStorage, error) {
	local := new(LocalStorage)
	local.LocalPath = path
	return local, os.MkdirAll(local.LocalPath, 0664)
}

// Put upload file
func (local *LocalStorage) Put(repo *StorageParse) error {
	srcFile, err := os.Open(repo.LocalPath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer srcFile.Close()

	desFile, err := os.Create(local.LocalPath + repo.RemotePath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer desFile.Close()

	_, err = io.Copy(desFile, srcFile)
	if err != nil {
		log.Error(err)
	}
	return err
}

// Get download file
func (local *LocalStorage) Get(repo *StorageParse) error {
	srcFile, err := os.Open(local.LocalPath + repo.RemotePath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer srcFile.Close()

	desFile, err := os.Create(repo.LocalPath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer desFile.Close()

	_, err = io.Copy(desFile, srcFile)
	if err != nil {
		log.Error(err)
	}
	return err
}

// Delete delete file
func (local *LocalStorage) Delete(repo *StorageParse) error {
	return os.Remove(local.LocalPath + repo.RemotePath)
}

// Share share file
func (local *LocalStorage) Share(repo *StorageParse) (string, error) {
	return "", nil
}
