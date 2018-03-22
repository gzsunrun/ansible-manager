package orm

import (
	log "github.com/astaxie/beego/logs"
	"time"
)

type Repository struct {
	ID      string                   `xorm:"repo_id" json:"repo_id"`
	Name    string                   `xorm:"repo_name" json:"repo_name"`
	Path    string                   `xorm:"repo_path" json:"repo_path"`
	Group   []map[string]interface{} `xorm:"repo_group" json:"repo_group"`
	Tag     []map[string]interface{} `xorm:"repo_tags" json:"repo_tags"`
	Vars    []Vars                   `xorm:"repo_vars" json:"repo_vars"`
	Note    string                   `xorm:"repo_notes" json:"repo_notes"`
	Desc    string                   `xorm:"repo_desc" json:"repo_desc"`
	Created time.Time                `xorm:"created" json:"created"`
}

type RepositoryInsert struct {
	ID      string                   `xorm:"repo_id" json:"repo_id"`
	Name    string                   `xorm:"repo_name" json:"repo_name"`
	Path    string                   `xorm:"repo_path" json:"repo_path"`
	Group   []map[string]interface{} `xorm:"repo_group" json:"repo_group"`
	Tag     []map[string]interface{} `xorm:"repo_tags" json:"repo_tags"`
	Vars    []Vars                   `xorm:"repo_vars" json:"repo_vars"`
	Note    string                   `xorm:"repo_notes" json:"repo_notes"`
	Desc    string                   `xorm:"repo_desc" json:"repo_desc"`
	Created time.Time                `xorm:"created" json:"-"`
}

type RepositoryList struct {
	ID      string                   `xorm:"repo_id" json:"repo_id"`
	Name    string                   `xorm:"repo_name" json:"repo_name"`
	Path    string                   `xorm:"repo_path" json:"-"`
	Group   []map[string]interface{} `xorm:"repo_group" json:"-"`
	Tag     []map[string]interface{} `xorm:"repo_tags" json:"-"`
	Vars    []Vars                   `xorm:"repo_vars" json:"-"`
	Note    string                   `xorm:"repo_notes" json:"-"`
	Desc    string                   `xorm:"repo_desc" json:"repo_desc"`
	Created time.Time                `xorm:"created" json:"created"`
}

type Vars struct {
	Name  string    `json:"vars_name"`
	Path  string    `json:"vars_path"`
	Value VarsValue `json:"vars_value"`
}

type VarsValue struct {
	Struct map[string]interface{} `json:"struct"`
	Vars   map[string]interface{} `json:"vars"`
}

func GetRepoByID(id string, repo interface{}) error {
	_, err := MysqlDB.Table("ansible_repository").Where("repo_id=?", id).Get(repo)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func FindRepos(repos interface{}) error {
	err := MysqlDB.Table("ansible_repository").Find(repos)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func CreateRepo(repo RepositoryInsert) error {
	_, err := MysqlDB.Table("ansible_repository").Insert(&repo)
	if err != nil {
		log.Error(err)
	}
	return err
}

func DelRepoByID(repoID string) error {
	_, err := MysqlDB.Exec("delete from ansible_repository where repo_id=?", repoID)
	if err != nil {
		log.Error(err)
	}
	return err
}
