package orm

import (
	"time"

	log "github.com/astaxie/beego/logs"
	
)

// Repository repo table
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

// RepositoryInsert repo insert
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

// RepositoryList repo output
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

// Vars repo vars
type Vars struct {
	Name  string    `json:"vars_name"`
	Path  string    `json:"vars_path"`
	Value VarsValue `json:"vars_value"`
}

// VarsValue vars value
type VarsValue struct {
	Struct map[string]interface{} `json:"struct"`
	Vars   map[string]interface{} `json:"vars"`
}

// GetRepoByID get repo by id
func GetRepoByID(id string, repo interface{}) error {
	_, err := MysqlDB.Table("ansible_repository").Where("repo_id=?", id).Get(repo)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// GetRepoByPath  repo exit by path
func GetRepoByPath(path string) (res bool){
	var repo Repository
	res, err := MysqlDB.Table("ansible_repository").Where("repo_path=?", path).Get(&repo)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

// GetRepoByName get repo by name
func GetRepoByName(name string)(*Repository,bool,error) {
	var repo Repository
	res, err := MysqlDB.Table("ansible_repository").Where("repo_name=?", name).Get(&repo)
	if err != nil {
		log.Error(err)
		return nil,false,err
	}
	return &repo,res,nil
}

// FindRepos find all repos
func FindRepos(repos interface{}) error {
	err := MysqlDB.Table("ansible_repository").Find(repos)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// CreateRepo insert repo into table
func CreateRepo(repo RepositoryInsert) error {
	_, err := MysqlDB.Table("ansible_repository").Insert(&repo)
	if err != nil {
		log.Error(err)
	}
	return err
}

// CreateRepos insert repo into table
func CreateRepos(repos []RepositoryInsert) error {
	session:=MysqlDB.NewSession()
	for _,v:=range repos{
		_,res,err:=GetRepoByName(v.Name)
		if err!=nil{
			return err
		}
		if res{
			log.Error(v.Name,"playbook is exit")
		}else{
			_,err:=session.Table("ansible_repository").Insert(v)
			if err!=nil{
				log.Error(err)
				return err
			}	
		}
	}
	err:=session.Commit()
	// _, err := MysqlDB.Table("ansible_repository").Insert(&repos)
	if err != nil {
		log.Error(err)
	}
	return err
}

// DelRepoByID delete repo
func DelRepoByID(repoID string) error {
	_, err := MysqlDB.Exec("delete from ansible_repository where repo_id=?", repoID)
	if err != nil {
		log.Error(err)
	}
	return err
}
