package storage

import(
	"io"
	"os"
	log "github.com/astaxie/beego/logs"
)

type LocalStorage struct{
	LocalPath string
}

func NewLocalStorage(path string)(*LocalStorage,error){
	local:=new(LocalStorage)
	local.LocalPath=path
	return local,os.MkdirAll(local.LocalPath, 0664)
}

func (local *LocalStorage)Put(repo *StorageParse)error{
	srcFile, err := os.Open(repo.LocalPath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer srcFile.Close()

	desFile, err := os.Create(local.LocalPath+repo.RemotePath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer desFile.Close()

	_,err=io.Copy(desFile, srcFile)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func (local *LocalStorage)Get(repo *StorageParse)error{
	srcFile, err := os.Open(local.LocalPath+repo.RemotePath)
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

	_,err=io.Copy(desFile, srcFile)
	if err!=nil{
		log.Error(err)
	}
	return err
}

func (local *LocalStorage)Delete(repo *StorageParse)error{
	return os.Remove(local.LocalPath+repo.RemotePath)
}

func (local *LocalStorage)Share(repo *StorageParse)(string,error){
	return "",nil
}