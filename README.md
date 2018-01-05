# ansible-manager
ansible-manager 实现的是类似ansible ui 的功能，让你的playbook脚本可视化，让你更方便管理你的playbook脚本,让部署更简单快捷。

## 如何使用

### 依赖

- ansible
- mariadb
- centos7

### 编译

- `go build -v`

### 启动

- 新建数据库`ansible_manager`,导入数据库`ansible_manager.sql`

- 配置`/etc/ansible_manager/ansible-manager.conf`
```
[ansible_manager]
port=8090
concurrent=5
work_path=/tmp/ansible-manager
# 数据库连接
mysql_url=root:123456@tcp(127.0.0.1:3306)/ansible?charset=utf8&loc=Local
# 登录token密钥
jwt_secret=XdCdkkffDM44DcDFSSF564bkDfffrcGMhfT0tyd3
# 是否使用s3作为脚本存储，是则需配置s3_endpoint、s3_key、s3_secret、bucket_name
s3=false
s3_endpoint=10.21.1.234:8080
s3_key=7YW507KJFWC0CYOXLSX0
s3_secret=XDCAx4atY96wELcSiwv9bkqCf0pCcWCGMhXToY56
bucket_name=ansible-playbook

```
- 启动`ansible-manager start`

- 访问`http://your ip:8090/ui/login.html` 默认账户:`admin`,密码:`123456`


## Contributing