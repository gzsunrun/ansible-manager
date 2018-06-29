package template

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	log "github.com/astaxie/beego/logs"
	"github.com/ghodss/yaml"
	"github.com/satori/go.uuid"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/tools/amcreate/template"
)

// Template the struct of AMfile.yml
type Template struct {
	Name 		string 	`json:"name"`
	Version 	string	`json:"version"`
	Desc		string	`json:"desc"`
	Inv			string	`json:"inventory"`
	AMDir 		string 	`json:"am_dir"`
	VarsDir 	string 	`json:"vars_dir"`
	Index		string	`json:"index"`
}

// ReadVars fetch vars from repo
func ReadVars(filePath string,remotePath string) ([]orm.RepositoryInsert,error) {
	repos:=make([]orm.RepositoryInsert,0)
	dir := filePath + "_dir"
	err := os.MkdirAll(dir, 0664)
	if err != nil {
		log.Error(err)
		return repos,err
	}
	cmd := exec.Command("tar", "zxvf", filePath, "-C", dir)
	err = cmd.Run()
	if err != nil {
		log.Error(err)
		return repos,err
	}

	err = template.Create(dir)
	if err != nil {
		log.Error(err)
		return repos,err
	}

	err = template.GzFile(dir,filePath)
	if err != nil {
		log.Error(err)
		return repos,err
	}
	
	tplByte, err := ioutil.ReadFile(dir + "/AMfile.yml")
	if err != nil {
		log.Error(err)
		return repos,err
	}

	var tpl []Template
	err = yaml.Unmarshal(tplByte, &tpl)
	if err != nil {
		log.Error(err)
		return repos,err
	}

	for _,t:=range tpl{
		tdir:=dir+"/"+t.AMDir
		var repo orm.RepositoryInsert
		repo.Name=t.Name+"-"+t.Version
		repo.Path=remotePath
		repo.Desc=t.Desc
		repo.ID=uuid.Must(uuid.NewV4()).String()
		groupY, err := ioutil.ReadFile(tdir + "/group.yml")
		if err != nil {
			log.Error(err)
			return repos,err
		}
		group, err := yaml.YAMLToJSON(groupY)
		if err != nil {
			log.Error(err)
			return repos,err
		}
		var groupMap []map[string]interface{}
		err = json.Unmarshal(group, &groupMap)
		if err != nil {
			log.Error(err)
			return repos,err
		}
		repo.Group = groupMap
		tagY, err := ioutil.ReadFile(tdir + "/tag.yml")
		if err != nil {
			log.Error(err)
			return repos,err
		}
		tag, err := yaml.YAMLToJSON(tagY)
		if err != nil {
			log.Error(err)
			return repos,err
		}
		var tagMap []map[string]interface{}
		err = json.Unmarshal(tag, &tagMap)
		if err != nil {
			log.Error(err)
			return repos,err
		}
		repo.Tag = tagMap
		repo.Note = "暂无简介"
		_, err = os.Stat(tdir + "/notes.md")
		if err == nil || os.IsExist(err) {
			notes, err := ioutil.ReadFile(tdir + "/notes.md")
			if err != nil {
				log.Error(err)
				return repos,err
			}
			repo.Note = string(notes)
		}

		varsPaths, err := getFilelist(tdir + "/vars")
		if err != nil {
			log.Error(err)
			return repos,err
		}
		vars := make([]orm.Vars, 0)
		for _, fpath := range varsPaths {
			vals := orm.Vars{}
			_, file := filepath.Split(fpath)
			vals.Name = strings.Replace(file, path.Ext(file), "", -1)
			vals.Path = t.VarsDir+"/"+file
			//vals.Path = strings.Replace(fpath, dir+"/", "", -1)
			val, err := ioutil.ReadFile(dir+"/"+vals.Path)
			if err != nil {
				log.Error(err)
				return repos,err
			}
			
			valJSON, err := yaml.YAMLToJSON(val)
			if err != nil {
				log.Error(err)
				return repos,err
			}
			vStruct, err := ioutil.ReadFile(fpath)
			if err != nil {
				log.Error(err)
				return repos,err
			}
			structJSON, err := yaml.YAMLToJSON(vStruct)
			if err != nil {
				log.Error(err)
				return repos,err
			}
			varsStr := `{"vars":` + string(valJSON) + `,"struct":` + string(structJSON) + `}`
			varsV := orm.VarsValue{}
			err = json.Unmarshal([]byte(varsStr), &varsV)
			if err != nil {
				log.Error(err)
				return repos,err
			}

			vals.Value = varsV
			vars = append(vars, vals)
		}
		repo.Vars = vars
		repos=append(repos,repo)
	}

	return repos,nil
}

// getFilelist walk a dir
func getFilelist(fpath string) ([]string, error) {
	paths := make([]string, 0)
	err := filepath.Walk(fpath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		log.Error(err)
	}
	return paths, err
}