## 如何使用

### 环境

- go 1.9
- ansible
- mariadb
- git
- centos7
- etcd3

### 编译
- `go build -v`

### 启动

- 安装mariadb数据库，并新建数据库 `ansible_manager`

- 导入数据库 `ansible_manager.sql`

```shell

mysql -u 用户名 -p密码  ansible_manager < ansible_manager.sql

```

- 静态文件public的文件放置 `/var/lib/ansible-manager/public/`

- 配置文件 `/etc/ansible-manager/ansible-manager.conf`，如下：
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

[git_storage]
enable=false

[file_log]
enable=true
log_dir=/opt/ansible-manager/log/

[etcd]
enable=true
endpoints=10.21.1.178:2379
```

**common** 节点配置

`port`： 服务端口

`concurrent`： 节点并发任务数

`work_path`:  工作路径

`master_enable`: 该节点是否作为master角色

`worker_enable`: 该节点是否作为worer角色

`uapi_enable`: 该节点是否开启web

`node_timeout`: 节点超时时间，超过该时间后无反应，认为该节点die

**mysql**  数据库配置

`mysql_url`: mysql地址,ip:port

`mysql_name`: mysql数据库名

`mysql_user`: mysql用户

`mysql_password`: mysql密码

**local_storage**   本地存储配置，用于存储脚本，local、s3、git选其一

`enable`: 是否应用，应用本地存储，只能够单机工作，无法建集群

`storage_dir`: 本地文件夹路径


**s3_storage**  S3存储配置，用于存储脚本，local、s3、git选其一

`enable`: 是否应用

`s3_endpoint`: S3地址

`s3_key`: S3 Key

`s3_secret`: S3 Secret

`bucket_name`: 桶名

**git_storage** git存储配置，用于存储脚本，local、s3、git选其一

`enable`: 是否应用

**file_log** 任务log信息保存在文件

`enable`: 是否应用

`log_dir`: 保存文件夹

**etcd** etcd配置

`enable`: 是否应用,不使用etcd将只能够单机工作，无法建集群

`endpoints`: etcd 地址 ip:port，多个使用`,`隔开



- 启动`ansible-manager start`

- 访问`http://your ip:8090/ui/login.html` 默认账户:`admin`,密码:`123456`
