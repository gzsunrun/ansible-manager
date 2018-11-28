package storage

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashwing/log"
)

// Git git repo
type Git struct {
}

// NewGit new git
func NewGit() *Git {
	return new(Git)
}

// Put upload file
func (g *Git) Put(repo *StorageParse) error {
	return nil
}

// Get download file
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

// Delete delete file
func (g *Git) Delete(repo *StorageParse) error {
	return nil
}

// Share share file
func (g *Git) Share(repo *StorageParse) (string, error) {
	return "", nil
}

// GetIO download file
func (local *Git) GetIO(repo *StorageParse) ([]byte, string, error) {
	return nil, "", nil
}
