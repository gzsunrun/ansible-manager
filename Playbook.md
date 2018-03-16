## 脚本规范

### 脚本压缩包

    1、压缩包格式 .tar.gz
    2、大小限制2G
    3、压缩包内容必须从根目录开始

### 文件目录
```
--vars
   --test1.yml
   --test2.yml
--ansible.cfg
--group.yml
--host
--index.yml
--tag.yml
--info.yml

```

#### vars 目录

放置全局参数配置文件，内容格式 ：yml 格式，文件命名：xxx.yml
    
#### hosts文件

内容如下：

```
{{range $index, $host:=.Hosts}}{{$host.HostName}} ansible_ssh_host={{$host.IP}} ansible_ssh_private_key_file=key-{{$host.IP}}
{{end}}

{{range $index, $group:=.Group}}[{{$group.Name}}]{{range $index, $host:=$group.Hosts}}
{{$host.HostName}} {{range $index, $attr:=$host.Attr}}{{$attr.Key}}={{$attr.Value}} {{end}}{{end}}

{{end}}
...... (固定内容)

```
     
**示例**：

```
host1 ansible_ssh_host=10.21.1.199 ansible_user=root ansible_ssh_pass=123456
host2 ansible_ssh_host=10.21.1.193 ansible_user=root ansible_ssh_pass=123456
host3 ansible_ssh_host=10.21.1.208 ansible_user=root ansible_ssh_pass=123456

[mariadb]
# mysql_master=yes (yes|no) 是否做为数据库集群引导节点,有且只有一个
host1 mysql_master=yes
host2
host3

[ntp:children]
mariadb

```

改变后：

```
{{range $index, $host:=.Hosts}}{{$host.HostName}} ansible_ssh_host={{$host.IP}} ansible_ssh_private_key_file=key-{{$host.IP}}
{{end}}

{{range $index, $group:=.Group}}[{{$group.Name}}]{{range $index, $host:=$group.Hosts}}
{{$host.HostName}} {{range $index, $attr:=$host.Attr}}{{$attr.Key}}={{$attr.Value}} {{end}}{{end}}

{{end}}


[ntp:children]
mariadb

```

#### group.yml 描述group信息
```yaml
- group_name: mariadb
  attr:
    - key: mysql_master
      type: bool
      default: "no"
```

`group_name`：group名字  
`attr`: 该group的group vars，数组  
`key`: group vars 的key  
`type`: group vars 的类型，bool，string  
`default`: 默认值 ,`bool 类型的必须是“yes","no",注意引号`  

#### index.yml 文件，脚本入口文件

#### ansible.cfg 文件，ansible配置文件

#### tag.yml 文件，playbook 运行tag
```yaml

- tag_name: 默认
  tag_value: ""

```

#### info.yml 文件，脚本信息
```yaml

repo_name: 脚本名字
repo_desc: 脚本描述

```

#### 其他自定义文件或文件夹
