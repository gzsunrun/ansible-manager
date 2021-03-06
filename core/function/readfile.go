package function

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashwing/log"
	"github.com/ghodss/yaml"
	"github.com/gzsunrun/ansible-manager/core/orm"
)

// ReadVars fetch vars from repo
func ReadVars(filePath string, repo *orm.RepositoryInsert) error {
	dir := filePath + "_dir"
	err := os.MkdirAll(dir, 0664)
	if err != nil {
		log.Error(err)
		return err
	}
	cmd := exec.Command("tar", "zxvf", filePath, "-C", dir)
	err = cmd.Run()
	if err != nil {
		log.Error(err)
		return err
	}
	groupY, err := ioutil.ReadFile(dir + "/group.yml")
	if err != nil {
		log.Error(err)
		return err
	}
	group, err := yaml.YAMLToJSON(groupY)
	if err != nil {
		log.Error(err)
		return err
	}
	var groupMap []map[string]interface{}
	err = json.Unmarshal(group, &groupMap)
	if err != nil {
		log.Error(err)
		return err
	}
	repo.Group = groupMap
	tagY, err := ioutil.ReadFile(dir + "/tag.yml")
	if err != nil {
		log.Error(err)
		return err
	}
	tag, err := yaml.YAMLToJSON(tagY)
	if err != nil {
		log.Error(err)
		return err
	}
	var tagMap []map[string]interface{}
	err = json.Unmarshal(tag, &tagMap)
	if err != nil {
		log.Error(err)
		return err
	}
	repo.Tag = tagMap
	repo.Note = "暂无简介"
	_, err = os.Stat(dir + "/notes.md")
	if err == nil || os.IsExist(err) {
		notes, err := ioutil.ReadFile(dir + "/notes.md")
		if err != nil {
			log.Error(err)
			return err
		}
		repo.Note = string(notes)
	}

	_, err = os.Stat(dir + "/info.yml")
	if err == nil || os.IsExist(err) {
		info, err := ioutil.ReadFile(dir + "/info.yml")
		if err != nil {
			log.Error(err)
			return err
		}

		err = yaml.Unmarshal(info, repo)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	varsPaths, err := getFilelist(dir + "/vars")
	if err != nil {
		log.Error(err)
		return err
	}
	vars := make([]orm.Vars, 0)
	for _, fpath := range varsPaths {
		if strings.HasSuffix(fpath, "_struct.yml") {
			continue
		}
		vals := orm.Vars{}
		val, err := ioutil.ReadFile(fpath)
		if err != nil {
			log.Error(err)
			return err
		}
		distr, file := filepath.Split(fpath)
		vals.Name = strings.Replace(file, path.Ext(file), "", -1)
		vals.Path = strings.Replace(fpath, dir+"/", "", -1)
		valJSON, err := yaml.YAMLToJSON(val)
		if err != nil {
			log.Error(err)
			return err
		}
		vStruct, err := ioutil.ReadFile(distr + "/" + vals.Name + "_struct.yml")
		if err != nil {
			log.Error(err)
			return err
		}
		structJSON, err := yaml.YAMLToJSON(vStruct)
		if err != nil {
			log.Error(err)
			return err
		}
		varsStr := `{"vars":` + string(valJSON) + `,"struct":` + string(structJSON) + `}`
		varsV := orm.VarsValue{}
		err = json.Unmarshal([]byte(varsStr), &varsV)
		if err != nil {
			log.Error(err)
			return err
		}

		vals.Value = varsV
		vars = append(vars, vals)
	}
	repo.Vars = vars
	return nil
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
