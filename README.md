# subscribe2clash

订阅转换clash配置，支持多订阅链接导入，多个链接以英文逗号分割

## 支持订阅类型

- v2ray
- ss
- ssr
- ssd

## 使用

### 二进制

- [release](https://github.com/whoisix/subscribe2clash/releases)下载对应的版本
- 解压后运行程序
- 访问http://localhost:8162/?sub_link=你的订阅链接

### 源码

- 安装Go 1.11+
- `go get github.com/whoisix/subscribe2clash`
- `export GO111MODULE=on`
- 编译 `make build`
- 启动 `./main`
- 访问http://localhost:8162/?sub_link=你的订阅链接

## 命令

- -h ：帮助
- -gc ：生成clash配置文件，不启动api服务
- -b ：clash配置基础文件，参考`config/clash/base_clash.yaml`
- -o ：配置文件名
- -l ：api服务监听地址
- -p ：api服务监听端口
- -origin ：acl规则获取地址。cn：国内镜像，github：github获取
- -proxy ：http代理

## 参考

- https://github.com/ne1llee/v2ray2clash

## 引用

- https://github.com/ACL4SSR/ACL4SSR

