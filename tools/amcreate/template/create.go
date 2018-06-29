package template

import (
	"strings"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"github.com/ghodss/yaml"
	"github.com/astaxie/beego/logs"
)

// Create create template
func Create(root string) error {
	data,err:=ioutil.ReadFile(root+"/AMfile.yml")
		if err!=nil{
			return err
		}
		var am []Template
		err=yaml.Unmarshal(data,&am)
		if err!=nil{
			return err
		}
		for _,v :=range am{

			// write vars
			varsPath,err:=getFilelist(root+"/"+v.VarsDir)
			if err!=nil{
				return err
			}
			amDir:=v.AmDir
			err=os.MkdirAll(root+"/"+amDir,0666)
			if err!=nil{
				return err
			}
			for _,p:=range varsPath{
				bpath:=strings.Replace(strings.Replace(p,"\\","/",-1),v.VarsDir,amDir+"/vars",1)
				err=os.MkdirAll(path.Dir(bpath),0666)
				if err!=nil{
					return err
				}
				data,err:=RefVars(p)
				if err!=nil{
					return err
				}
				ioutil.WriteFile(bpath,data,0666)
				if err!=nil{
					return err
				}
			}

			// write group.yml
			group,inv,err:=GetGroup(root+"/"+v.Inv)
			if err!=nil{
				return err
			}
			inv=inventory+inv
			err=ioutil.WriteFile(root+"/"+amDir+"/hosts",[]byte(inv),0666)
			if err!=nil{
				return err
			}

			err=ioutil.WriteFile(root+"/"+amDir+"/group.yml",[]byte(group),0666)
			if err!=nil{
				return err
			}

			// write ansible.cfg
			cfg:=""
			err=readLine(root+"/"+path.Dir(v.Index)+"/ansible.cfg",func(p string){
				
				// inventory inventory path
				if strings.Contains(p,"inventory"){
					cfg+="inventory = hosts\n"
					return
				}

				// roles path
				if strings.Contains(p,"roles_path"){
					
					ps:=strings.Split(p,"=")
					rolePath:=strings.TrimSpace(ps[1])
					for _,in:=range strings.Split(v.AmDir,"/"){
						if in!=""{
							rolePath ="../"+rolePath
						}
					}
					cfg+="roles_path = "+rolePath+"\n"
					return
				}
				cfg+=p+"\n"
			})

			err = ioutil.WriteFile(root+"/"+amDir+"/ansible.cfg",[]byte(cfg),0666)
			if err!=nil{
				return err
			}

			// write tag.yml
			err = ioutil.WriteFile(root+"/"+amDir+"/tag.yml",[]byte(tag),0666)
			if err!=nil{
				logs.Error(err)
				return err
			}

			// write notes.md
			notes := []byte("## "+v.Name)
			if v.Notes != ""{
				notes,err = ioutil.ReadFile(root+"/"+v.Notes)
				if err!=nil{
					logs.Error(err)
					return err
				}
			}
			err = ioutil.WriteFile(root+"/"+amDir+"/notes.md",notes,0666)
			if err!=nil{
				logs.Error(err)
				return err
			}

			// write index.yml
			for _,in:=range strings.Split(v.AmDir,"/"){
				if in!=""{
					v.Index="../"+v.Index
				}
			}
			err = ioutil.WriteFile(root+"/"+amDir+"/index.yml",[]byte(`- include: "`+v.Index+`"`),0666)
			if err!=nil{
				logs.Error(err)
				return err
			}
		}
		return nil
}

// Template the struct of AMfile.yml
type Template struct {
	Name 		string 	`json:"name"`
	Version 	string	`json:"version"`
	Desc		string	`json:"desc"`
	Notes		string	`json:"notes"`
	Inv			string	`json:"inventory"`
	AmDir 		string 	`json:"am_dir"`
	VarsDir 	string 	`json:"vars_dir"`
	Index		string	`json:"index"`
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
		logs.Error(err)
	}
	return paths, err
}

var tag =`- tag_name: 默认
  tag_value: ""
`
var inventory =`{{range $index, $host:=.Hosts}}{{$host.HostName}} ansible_ssh_host={{$host.IP}} {{if ne $host.Key ""}}ansible_ssh_private_key_file=key-{{$host.IP}}{{else}}ansible_user={{$host.User}} ansible_ssh_pass={{$host.Password}}{{end}} 
{{end}}
		
{{range $index, $group:=.Group}}[{{$group.Name}}]{{range $index, $host:=$group.Hosts}}
{{$host.HostName}} {{range $index, $attr:=$host.Attr}}{{$attr.Key}}={{$attr.Value}} {{end}}{{end}}
		
{{end}}
`


// GzFile gz root to dest
func GzFile(root ,dest string) error{
	files :=make([]*os.File,0)
	dir, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, fi := range dir {
		f,err:=os.Open(root+"/"+fi.Name())
		if err!=nil{
			logs.Error(err)
			return err
		}
		files=append(files,f)
	}
		
	err = Compress(files,dest)
	if err!=nil{
		logs.Error(err)
		return err
	}
	return nil
}