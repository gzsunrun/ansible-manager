package helm

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/hashwing/log"

	"github.com/ghodss/yaml"
	"github.com/gzsunrun/ansible-manager/core/function"
	"github.com/gzsunrun/ansible-manager/core/orm"
	"github.com/gzsunrun/ansible-manager/tools/amcreate/template"
)

// deCompressChart DeCompressChart to dest
func deCompressChart(tarFile, dest string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		err = os.MkdirAll(path.Dir(filename), 0664)
		if err != nil {
			return err
		}
		fmt.Println(filename)
		fw, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer fw.Close()
		_, err = io.Copy(fw, tr)
		if err != nil {
			return err
		}
	}

	return nil
}

type ChartInfo struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Desc       string `json:"description"`
	AppVersion string `json:"appVersion"`
	Icon       string `json:"icon"`
}

// ReadChart
func ReadChart(tarFile, remotePath string) ([]orm.RepositoryInsert, error) {
	dir := tarFile + "_dir/"
	err := deCompressChart(tarFile, dir)
	if err != nil {
		return nil, err
	}
	fileInfo, err := ioutil.ReadDir(dir)
	if len(fileInfo) == 0 {
		return nil, errors.New(dir + " have no file")
	}
	if fileInfo[0].IsDir() {
		dir += fileInfo[0].Name() + "/"
	}

	valuesByte, err := ioutil.ReadFile(dir + "values.yaml")
	if err != nil {
		return nil, err
	}

	var values map[string]interface{}
	err = yaml.Unmarshal(valuesByte, values)
	if err != nil {
		return nil, err
	}
	infoByte, err := ioutil.ReadFile(dir + "Chart.yaml")
	if err != nil {
		return nil, err
	}
	info := ChartInfo{}
	err = yaml.Unmarshal(infoByte, &info)
	if err != nil {
		return nil, err
	}

	if info.Icon != "" {
		err := os.MkdirAll(tarFile+"_dir/logo", 0664)
		if err != nil {
			return nil, err
		}
		download(info.Icon, tarFile+"_dir/logo/"+remotePath)
	}

	valuesStruct, err := template.RefVarsValue(valuesByte)
	if err != nil {
		return nil, err
	}

	varValue := orm.VarsValue{
		Struct: valuesStruct,
		Vars:   values,
	}
	vars := orm.Vars{
		Name:  "default",
		Path:  "values.yaml",
		Value: varValue,
	}
	log.Debug(info.Desc)

	return []orm.RepositoryInsert{orm.RepositoryInsert{
		ID:      function.NewUuidString(),
		Name:    info.Name,
		Version: info.AppVersion + "@" + info.Version,
		Type:    "helm",
		Path:    remotePath,
		Vars:    []orm.Vars{vars},
		Desc:    info.Desc,
	}}, nil

}

func download(UrlStr, name string) error {

	res, err := http.Get(UrlStr)
	if err != nil {
		return err
	}
	f, err := os.Create(name + path.Ext(UrlStr))
	if err != nil {
		return err
	}

	_, err = io.Copy(f, res.Body)
	return err

}
