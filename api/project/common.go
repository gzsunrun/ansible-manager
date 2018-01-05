package project

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/astaxie/beego/logs"
)

func JsonWrite(w http.ResponseWriter, code int, msg interface{}) {
	w.WriteHeader(code)
	data, _ := json.Marshal(&msg)
	w.Write(data)
}

type varsStruct struct {
	group string
	tag   string
	vars  []map[string]string
}

func readVars(filePath string) (*varsStruct, error) {
	vst := &varsStruct{}
	dir := filePath + "_dir"
	err := os.MkdirAll(dir, 0664)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	defer os.RemoveAll(dir)
	cmd := exec.Command("tar", "zxvf", filePath, "-C", dir)
	err = cmd.Run()
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	groupY, err := ioutil.ReadFile(dir + "/group.yml")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	group, err := yaml.YAMLToJSON(groupY)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	vst.group = string(group)
	tagY, err := ioutil.ReadFile(dir + "/tag.yml")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	tag, err := yaml.YAMLToJSON(tagY)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	vst.tag = string(tag)
	varsPaths, err := getFilelist(dir + "/vars")
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	vars := make([]map[string]string, 0)
	for _, fpath := range varsPaths {
		vals := make(map[string]string)
		val, err := ioutil.ReadFile(fpath)
		if err != nil {
			logs.Error(err)
			return nil, err
		}
		_, file := filepath.Split(fpath)
		vals["name"] = strings.Replace(file, path.Ext(file), "", -1)
		vals["path"] = strings.Replace(fpath, dir+"/", "", -1)
		vals["value"] = string(val)
		vars = append(vars, vals)
	}
	vst.vars = vars
	return vst, nil
}

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
		logs.Error(err)
	}
	return paths, err
}
