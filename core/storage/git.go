package storage

import (
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/astaxie/beego/logs"
)

type Git struct {
}

func NewGit() *Git {
	return new(Git)
}

func (g *Git) Put(repo *StorageParse) error {
	return nil
}

func (g *Git) Get(repo *StorageParse) error {
	localDir := repo.LocalPath + "_git"
	cmd := exec.Command("git", "clone", repo.RemotePath, localDir)
	cmd.Dir = filepath.Dir(repo.LocalPath)
	err := cmd.Run()
	if err != nil {
		log.Error("clone:", err)
		return err
	}
	defer os.RemoveAll(localDir)

	cmd = exec.Command("tar", "zcvf", repo.LocalPath, "./")
	cmd.Dir = localDir
	return cmd.Run()
}

func (g *Git) Delete(repo *StorageParse) error {
	return nil
}

func (g *Git) Share(repo *StorageParse) (string, error) {
	return "", nil
}
