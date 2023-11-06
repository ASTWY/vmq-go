# 部署

## 1. 下载可执行文件

从 [GitHub Release](https://github.com/astwy/vmq-go/releases) 下载最新的可执行文件压缩包。解压后可得到可执行文件  
解包后的目录结构如下:

```bash
.
├── config.example.yaml # 配置文件样例
├── logs # 日志目录 (运行后自动生成)
├── web # 静态文件目录
└── vmq-go # 可执行文件 根据系统不同名称可能不同
```

## 2. 配置文件

需要将解压后的 `config.example.yaml` 文件重命名为 `config.yaml`  
配置项详情请参考 [配置文件](/config/)

## 3. 运行

执行可执行文件即可运行

## 4. 访问

默认监听地址为 `服务器IP:8080`

## 5. 后台管理

后台登录地址为 `服务器IP:8080/#/login`  
默认管理员账号为 `admin` 密码为 `admin`

## 6. 配置 VMQ

登录到后台管理后，点击左侧菜单栏的 `设置`，根据提示配置 VMQ
