package storage

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/hashwing/log"
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

// GetIO download file
func (local *LocalStorage) GetIO(repo *StorageParse) ([]byte, string, error) {
	infos, err := ioutil.ReadDir(local.LocalPath + "/logos")
	if err != nil {
		return nil, "", err
	}
	for _, info := range infos {
		if strings.HasPrefix(info.Name(), repo.LocalPath) {
			data, err := ioutil.ReadFile(local.LocalPath + "/logos/" + info.Name())
			return data, path.Ext(info.Name()), err
		}

	}
	return nil, "", nil
}
