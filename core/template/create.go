package template

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/ghodss/yaml"
	"github.com/gzsunrun/ansible-manager/core/function"
	"github.com/gzsunrun/ansible-manager/core/orm"
)

// PlaybookParse playbook parse struct
type PlaybookParse struct {
	Hosts []orm.HostsList `json:"hosts"`
	Group []orm.Group     `json:"group"`
	Vars  []orm.Vars      `json:"vars"`
}

// ReadAMfile read am file
func ReadAMfile(dir string)([]Template,error){
	var tpls []Template
	tplByte, err := ioutil.ReadFile(dir + "/AMfile.yml")
	if err != nil {
		return tpls,err
	}

	err = yaml.Unmarshal(tplByte, &tpls)
	if err != nil {
		return tpls,err
	}
	return tpls,err
}

// InstallVars create group vars
func InstallVars(t *orm.Task,dir string) (workPath string ,err error) {
	tplByte, err := ioutil.ReadFile(dir + "/AMfile.yml")
	if err != nil {
		return
	}

	var tpls []Template
	err = yaml.Unmarshal(tplByte, &tpls)
	if err != nil {
		return
	}

	var repo orm.Repository
	err = orm.GetRepoByID(t.RepoID,&repo)
	if err != nil {
		return
	}
	var temp Template
	for _,tpl:=range tpls{
		if tpl.Name+"-"+tpl.Version==repo.Name{
			temp=tpl
			workPath = tpl.AMDir
			break
		}
	}
	
	var hosts []orm.HostsList
	err = orm.FindHostFromProject(t.ProjectID, &hosts)
	if err != nil {
		return
	}
	for i, g := range t.Group {
		for j, h := range g.Hosts {
			for _, hh := range hosts {
				if hh.ID == h.HostUUID {
					t.Group[i].Hosts[j].IP = hh.IP
					t.Group[i].Hosts[j].HostName = hh.HostName
				}
			}
		}
	}
	for i, val := range hosts {
		if val.Key != "" && val.Password != "" {
			if function.AuthKey(val) != "success" {
				hosts[i].Key = ""
			}
		}
		if val.HostName == "" {
			hosts[i].HostName = "host" + strconv.Itoa(i)
		}
		if val.Key != "" {
			err = ioutil.WriteFile(dir+"/"+temp.AMDir+"/key-"+val.IP, []byte(val.Key), 0600) 
			if err != nil {
				return
			}
		}
	}
	playbookParse := PlaybookParse{
		Hosts: hosts,
		Group: t.Group,
		Vars:  t.Vars,
	}

	for _, val := range t.Vars {
		v, err1 := json.Marshal(val.Value.Vars)
		if err1 != nil {
			err=err1
			return
		}
		err = ioutil.WriteFile(dir+"/"+val.Path, v, 0600)
		if err != nil {
			return
		}
	}

	tmpl, err := template.ParseFiles(dir + "/" + temp.AMDir + "/hosts")
	if err != nil {
		return
	}
	fd, err := os.OpenFile(dir+"/"+temp.AMDir+"/hosts", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return
	}
	err = tmpl.Execute(fd, playbookParse)
	if err != nil {
		return
	}
	return
}