package template

import (
	"strings"
	"io/ioutil"
	"github.com/ghodss/yaml"
	"github.com/buger/jsonparser"
	"github.com/astaxie/beego/logs"
)

// VarsStruct var struct
type VarsStruct struct{
	Head *map[string]interface{}
	Stack *DataStack
}
// DataStack data stack
type  DataStack struct {
	Top 	int
	Data 	[]*map[string]interface{}
}

//NewStack new a stack
func NewStack(size int)*DataStack{
	return &DataStack{
		Top:0,
		Data:make([]*map[string]interface{},size),
	}
}

// Push push item into stack
func (ds *DataStack)Push(item *map[string]interface{}){
	ds.Data[ds.Top]=item
	ds.Top++
}

// Pop pop item from stack
func (ds *DataStack)Pop()*map[string]interface{}{
	data:=ds.Data[ds.Top-1]
	ds.Top--
	return data
}

// Get get stack top item
func (ds *DataStack)Get()*map[string]interface{}{
	data:=ds.Data[ds.Top-1]
	return data
}

// Post global stack
var Post = VarsStruct{}

// RefVars ref var to structs
func RefVars(path string)([]byte,error){
	y,err:=ioutil.ReadFile(path)
	if err!=nil{
		logs.Error(err)
		return nil,err
	}
	ystr:=strings.Replace(string(y),"#@","",-1)
	ystr=strings.Replace(ystr,"[]","",-1)
	j,err:=yaml.YAMLToJSON([]byte(ystr))
	if err!=nil{
		logs.Error(err)
		return nil,err
	}
	stack:=NewStack(500)
	data:=make(map[string]interface{})
	Post.Head=&data
	Post.Stack=stack
	stack.Push(Post.Head)
	err=jsonparser.ObjectEach(j,refCallback)
	//fmt.Println(string(d))
	dd,err:=yaml.Marshal(&data)
	if err!=nil{
		logs.Error(err)
		return nil,err
	}
	//fmt.Println(string(dd))
	return dd,nil
}

func refCallback(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
	//fmt.Println(string(key))
	if dataType.String()=="object"{
		defer Post.Stack.Pop()
		p :=make(map[string]interface{})
		cur:=map[string]interface{}{
			"type":"object",
			"properties":&p,
		}
		(*Post.Stack.Get())[string(key)]=cur
		Post.Stack.Push(&p)
		err:=jsonparser.ObjectEach(value,refCallback)
		return err
	}
	if dataType.String()=="array"{
		next :=make(map[string]interface{})
		cur:=map[string]interface{}{
			"type":"array",
			"items":&next,
		}
		var err error
		value, dataType,offset,err =jsonparser.Get(value,"[0]")
		if err!=nil{
			logs.Error(err,",maybe this is a null array, key:",string(key),"item type:",dataType.String())
		}

		if dataType.String()=="object"{
			defer Post.Stack.Pop()
			p :=make(map[string]interface{})
			next =map[string]interface{}{
				"type":"object",
				"properties":&p,
			}
			(*Post.Stack.Get())[string(key)]=cur
			Post.Stack.Push(&p)
			err =jsonparser.ObjectEach(value,refCallback)
			if err!=nil{
				logs.Error(err)
			}
			return err
		}
		next=map[string]interface{}{
			"type":dataType.String(),
		}
		(*Post.Stack.Get())[string(key)]=cur
		return nil
		
	}

	cur:=map[string]interface{}{
		"type":dataType.String(),
	}
	(*Post.Stack.Get())[string(key)]=cur
	return nil
}
