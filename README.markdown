# V 免签 Go 版

## [文档](https://astwy.github.io/vmq-go/)

## 介绍

参考 [V 免签](https://github.com/szvone/vmqphp) PHP 版本，使用 Go 语言实现的 V 免签服务端。

## 特点

- 部署简单，无需安装 PHP 环境，设置伪静态等，下载文件，修改配置，运行即可。
- 支持 Docker Compose 快速部署
- 兼容原监控端接口
- 后台管理界面更加美观，操作更加方便
- 增设收款日志，方便查询收款记录
- 增设邮件通知功能
  - 监控端掉线通知
  - 收款通知
  - 收款异常通知

## 使用

### 1. 下载

在 [Release](/astwy/vmq-go/releases) 页面或 [Actions](/astwy/vmq-go/actions) 页面下载对应系统的可执行文件。

### 2. 配置

将下载的压缩包解压，修改**config.example.yam**文件，配置数据库等信息。然后将文件重命名**config.yaml**

### 3. 运行

根据不同的系统，运行对应的可执行文件即可。

## TODO

- [ ] 优化日志
- [ ] 后台安全入口
- [ ] 使用文档
