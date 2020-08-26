<h1 align="center">
  <br>subscribe2clash<br>
</h1>


<h4 align="center">Clash规则配置转换</h4>

<p align="center">
  <a href="https://github.com/whoisix/subscribe2clash/actions">
    <img src="https://img.shields.io/github/workflow/status/whoisix/subscribe2clash/Go" alt="Github Actions">
  </a>
  <a href="https://goreportcard.com/report/github.com/whoisix/subscribe2clash">
    <img src="https://goreportcard.com/badge/github.com/whoisix/subscribe2clash">
  </a>
  <a href="https://github.com/whoisix/subscribe2clash/releases">
    <img src="https://img.shields.io/github/release/whoisix/subscribe2clash/all.svg">
  </a>
</p>


## 简介

Clash规则配置转换，自动获取[ACL4SSR](https://github.com/ACL4SSR/ACL4SSR)路由规则。  

支持v2ray\trojan\ss\ssr\ssd订阅转换。  

支持多订阅一起转换，多个订阅连接用英文逗号隔开。

## 启动服务

### 二进制

- [release](https://github.com/whoisix/subscribe2clash/releases)下载对应的版本
- 解压后执行`./subscribe2clash -b 你的基础配置文件 -origin cn`
- 访问http://localhost:8162/?sub_link=你的订阅链接

### 源码

- 安装Go 1.11+
- `go get github.com/whoisix/subscribe2clash`
- `export GO111MODULE=on`
- 编译 `make build`
- 启动 `./main`
- 访问http://localhost:8162/?sub_link=你的订阅链接

## 命令

- 如果只想生成clash配置文件（没有节点数据），不启用api服务，可使用命令

  ```
  ./main -gc
  ```

- 指定自定义基础配置文件，可在里面添加自定义的路由规则，程序将按照这个文件写入路由信息，可以参考[config/clash/base_clash.yaml](https://github.com/whoisix/subscribe2clash/blob/master/config/clash/base_clash.yaml)，`%s`将被程序替换为ACL的路由规则。

  ```
  ./main -b ./yourfile.yaml
  ```

- 指定输出的配置文件。默认情况下配置文件会输出为`./config/clash/acl.yaml`，可以通过以下命令来重新指定。

  ```
  ./main -o ./yourconfig.yaml
  ```

- 获取ACL规则的源地址。cn：国内镜像（更新可能没有github及时），github：github获取。默认从github获取。

  ```
  ./main -origin github
  ./main -origin cn
  ```

- 启用http代理。由于网络原因，ACL的github源可能连接不上，但又不想使用镜像时，你可能需要配合代理一起食用。

  ```
  ./main -proxy http://127.0.0.1:7890
  ```

- 指定api服务监听端口，默认监听8162端口。

  ```
  ./main -p 8162
  ```

- 指定更新规则频率，单位小时，默认每6小时拉取一次。

  ```
  ./main -t 6
  ```

  

## 参考

- https://github.com/ne1llee/v2ray2clash

## 引用

- https://github.com/ACL4SSR/ACL4SSR

## 测试地址
http://47.106.211.213:8162/?sub_link=yourlink
