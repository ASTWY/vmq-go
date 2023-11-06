# 配置文件

## 配置文件路径

配置文件默认路径为 `./config.yaml`

## 配置文件说明

[样例配置文件](https://github.com/ASTWY/vmq-go/blob/main/config.example.yaml)中包含了所有可用的配置项，以及其默认值。  

```yaml
host: 0.0.0.0 # 监听地址
port: 8080 # 监听端口
log: # 日志配置
  level: info # 日志等级
  path: ./logs # 日志路径
jwt: # JWT配置
  secret: hfjtSdYL6j!Ts$uD # JWT密钥 (必填, 请自行修改)
  expire: 3600 # JWT过期时间 (单位: 秒)
db: # 数据库配置 (目前仅支持MySQL)
  host: localhost # 数据库地址
  port: 3306 # 数据库端口
  user: test # 数据库用户名
  password: test # 数据库密码
  name: test # 数据库名
```
