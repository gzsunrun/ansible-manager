# ansible-manager [![Build Status](https://travis-ci.org/gzsunrun/ansible-manager.svg?branch=v2.0)](https://travis-ci.org/gzsunrun/ansible-manager)
ansible-manager 实现的是类似ansible ui 的功能，让你的playbook脚本可视化，让你更方便管理你的playbook脚本,让部署更简单快捷。

## 如何使用

### 环境

- go 1.6+
- ansible
- mariadb
- centos7
- etcd3

### 编译
- `go build -v`

### 启动

- 新建数据库`ansible_manager`,导入数据库表结构`ansible_manager.sql`

- 静态文件public的文件放置 `/var/lib/ansible-manager/public/`

- 配置文件 `/etc/ansible-manager/ansible-manager.conf`
```
[common]
port=8090
concurrent=5
work_path=/tmp/ansible-manager
master_enable=true
worker_enable=true
uapi_enable=true
node_timeout=10


[mysql]
mysql_url=10.21.1.178:3306
mysql_name=playbook
mysql_user=root
mysql_password=123456

[local_storage]
enable=false
storage_dir=/opt/ansible-manager/repo/


[s3_storage]
enable=true
s3_endpoint=10.21.1.234:8080
s3_key=V7TRJJ6RMA8SWKDZEW6X
s3_secret=cNaO54RUWpDypOIuLGDJAIrqeq1vXeQ2QvZzLHMV
bucket_name=ansible-playbook

[file_log]
enable=true
log_dir=/opt/ansible-manager/log/

[etcd]
enable=true
endpoints=10.21.1.178:2379
```

- 启动`ansible-manager start`

- 访问`http://your ip:8090/ui/login.html` 默认账户:`admin`,密码:`123456`


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
host1 ansible_ssh_host=10.21.1.199 ansible_user=root ansible_ssh_pass=sunrunvas
host2 ansible_ssh_host=10.21.1.193 ansible_user=root ansible_ssh_pass=sunrunvas
host3 ansible_ssh_host=10.21.1.208 ansible_user=root ansible_ssh_pass=sunrunvas

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

#### 其他自定义文件或文件夹


## Contributing