package main

import (
	"fmt"
	"strings"
	"io"
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"github.com/ghodss/yaml"
	"github.com/astaxie/beego/logs"
	"github.com/gzsunrun/ansible-manager/tools/amcreate/template"
)

func main(){
	if len(os.Args)!=2{
		fmt.Print("error,please try:\n\namcreate build : create a template\namcreate gz : create .tar.gz\n\n")
		return 
	}
	root:=path.Dir(os.Args[0])

	if os.Args[1]=="build"{
		data,err:=ioutil.ReadFile(root+"/AMfile.yml")
		if err!=nil{
			logs.Error(err)
			return
		}
		var am []Template
		err=yaml.Unmarshal(data,&am)
		if err!=nil{
			logs.Error(err)
			return
		}
		for _,v :=range am{
			varsPath,err:=getFilelist(root+"/"+v.VarsDir)
			if err!=nil{
				logs.Error(err)
				return
			}
			amDir:=v.AmDir
			err=os.MkdirAll(root+"/"+amDir,0666)
			if err!=nil{
				logs.Error(err)
				return
			}
			logs.Info(varsPath,"\n")
			for _,p:=range varsPath{
				logs.Info(p,"\n")
				bpath:=strings.Replace(strings.Replace(p,"\\","/",-1),v.VarsDir,amDir+"/vars",1)
				logs.Info(bpath)
				err=os.MkdirAll(path.Dir(bpath),0666)
				if err!=nil{
					logs.Error(err)
					return
				}
				data,err:=template.RefVars(p)
				if err!=nil{
					logs.Error(err)
					return
				}
				ioutil.WriteFile(bpath,data,0666)
				if err!=nil{
					logs.Error(err)
					return
				}
			}
			group,inv,err:=template.GetGroup(root+"/"+v.Inv)
			if err!=nil{
				return
			}
			inv=inventory+inv
			err=ioutil.WriteFile(root+"/"+amDir+"/hosts",[]byte(inv),0666)
			if err!=nil{
				logs.Error(err)
				return
			}
			err=ioutil.WriteFile(root+"/"+amDir+"/group.yml",[]byte(group),0666)
			if err!=nil{
				logs.Error(err)
				return
			}
			cfg:=""
			err=readLine(root+"/"+path.Dir(v.Index)+"/ansible.cfg",func(p string){
				if strings.Contains(p,"inventory"){
					cfg+="inventory = hosts\n"
					return
				}
				cfg+=p+"\n"
			})
			if err!=nil{
				logs.Error(err)
				return
			}
			err = ioutil.WriteFile(root+"/"+amDir+"/ansible.cfg",[]byte(cfg),0666)
			if err!=nil{
				logs.Error(err)
				return
			}
			err = ioutil.WriteFile(root+"/"+amDir+"/tag.yml",[]byte(tag),0666)
			if err!=nil{
				logs.Error(err)
				return
			}
			err = ioutil.WriteFile(root+"/"+amDir+"/notes.md",[]byte("## "+v.Name),0666)
			if err!=nil{
				logs.Error(err)
				return
			}
			for _,in:=range strings.Split(v.AmDir,"/"){
				if in!=""{
					v.Index="../"+v.Index
				}
			}
			err = ioutil.WriteFile(root+"/"+amDir+"/index.yml",[]byte(`- include: "`+v.Index+`"`),0666)
			if err!=nil{
				logs.Error(err)
				return
			}
		}
	}
	
	if os.Args[1]=="gz"{
		files :=make([]*os.File,0)
		dir, err := ioutil.ReadDir(root)
		if err != nil {
			return
		}
		for _, fi := range dir {
			f,err:=os.Open(root+"/"+fi.Name())
			if err!=nil{
				logs.Error(err)
				return
			}
			files=append(files,f)
		}
		
		err = template.Compress(files,"test.tar.gz")
		if err!=nil{
			logs.Error(err)
			return
		}
	}
}

// Template the struct of AMfile.yml
type Template struct {
	Name 		string 	`json:"name"`
	Version 	string	`json:"version"`
	Desc		string	`json:"desc"`
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

func readLine(fileName string, handler func(string)) error {  
    f, err := os.Open(fileName)  
    if err != nil {  
        return err  
    }  
    buf := bufio.NewReader(f)  
    for {  
        line, err := buf.ReadString('\n')  
        line = strings.TrimSpace(line)  
        handler(line)  
        if err != nil {  
            if err == io.EOF {  
                return nil  
            }  
            return err  
        }  
    }  
    return nil  
} 
