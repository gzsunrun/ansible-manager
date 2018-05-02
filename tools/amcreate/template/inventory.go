package template

import (
	"strconv"
	"bufio"  
    "fmt"  
    "io"  
	"os"  
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
)

// Group ansible inventory group
type Group struct{
	Name	string 		`json:"group_name"`
	Attr	[]GroupAttr	`json:"attr"`
}

// GroupAttr group attr
type GroupAttr struct{
	Key		string	`json:"key"`
	Default	string	`json:"default"`
	Type	string	`json:"type"`
}

// GetGroup get group struct by inventory file
func GetGroup(path string)(string,string,error){
	groupMap:=make(map[string]*Group)
	groupName:=""
	staticContent:=""
	child:=false
	err:=readLine(path,func (s string){
		if strings.Contains(s,"#"){
			return 
		}
		res1,err:=regexp.MatchString(`\[[a-z]*\]`,s)
		if err!=nil{
			return
		}
		if res1{
			child=false
			digitsRegexp := regexp.MustCompile(`[a-z]+`)
			name:=digitsRegexp.FindStringSubmatch(s)
			groupName = name[0]
			if groupMap[groupName]==nil{
				groupMap[groupName]=new(Group)
				groupMap[groupName].Name=groupName
				groupMap[groupName].Attr=make([]GroupAttr,0)
			}
			return
		}
		res2,err:=regexp.MatchString(`\[[a-z]*:vars\]`,s)
		if err!=nil{
			return
		}
		if res2{
			child=false
			digitsRegexp := regexp.MustCompile(`[a-zA-Z]+:[a-zA-Z]+`)
			name:=digitsRegexp.FindStringSubmatch(s)
			groupName = strings.Split(name[0],":")[0]
			if groupMap[groupName]==nil{
				groupMap[groupName]=new(Group)
				groupMap[groupName].Name=groupName
				groupMap[groupName].Attr=make([]GroupAttr,0)
			}
			return
		}
		res3,err:=regexp.MatchString(`\[[a-z]*:children\]`,s)
		if err!=nil{
			return
		}
		if res3{
			staticContent+=s+"\n"
			child=true
			return
		}
		if groupName == ""{
			return
		}
		if child{
			staticContent+=s+"\n"
			return
		}
		fileds:=strings.Split(s," ")
		for _,filed:=range fileds{
			if filed==""{
				continue
			}
			res,err:=regexp.MatchString(`.=.`,filed)
			if err!=nil{
				return
			}
			if res{
				attr:=strings.Split(filed,"=")
				if len(attr)==2{
					if !strings.HasPrefix(attr[0],"ansible_"){
						gattr:=GroupAttr{}
						for _,g:=range groupMap[groupName].Attr{
							if g.Key==attr[0]{
								return
							}
						}
						gattr.Key=attr[0]
						gattr.Default=attr[1]
						gattr.Type="string"
						if attr[1]=="yes"||attr[1]=="no"||attr[1]=="true"||attr[1]=="false"{
							gattr.Type="bool"
						}
						if _,err:=strconv.Atoi(attr[0]);err==nil{
							gattr.Type="number"
						}
						groupMap[groupName].Attr=append(groupMap[groupName].Attr,gattr)
					}
				}
					
			}
		}
		
	})
	if err!=nil{
		fmt.Println(err)
		return "","",err
	}
	groups:=make([]Group,0)
	for _,v:=range groupMap{
		groups=append(groups,*v)
	}
	data,err:=yaml.Marshal(groups)
	if err!=nil{
		fmt.Println(err)
		return "","",err 
	}
	//fmt.Println(string(data))
	//fmt.Println(staticContent)
	return string(data),staticContent,nil
}


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