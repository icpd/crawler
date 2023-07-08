<h1 align="center">
  <img src="https://github.com/icpd/subscribe2clash/raw/master/subscribe2clash.png" alt="subscribe2clash" width="350">
</h1>


<p align="center">
  <a href="https://goreportcard.com/report/github.com/icpd/subscribe2clash">
    <img src="https://goreportcard.com/badge/github.com/icpd/subscribe2clash">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/Dreamacro/clash?style=flat-square">
  <a href="https://github.com/icpd/subscribe2clash/releases">
    <img src="https://img.shields.io/github/release/icpd/subscribe2clash/all.svg">
  </a>
</p>


## 简介

Clash配置转换，默认自动获取[ACL4SSR](https://github.com/ACL4SSR/ACL4SSR)路由规则。  

支持v2ray\trojan\ss\ssr\ssd订阅转换。**暂不支持单节点转换**  

支持多订阅一起转换，多个订阅连接用英文逗号隔开。

### 改动点
在原版基础上扩展了配置单个ss服务器的功能
新增页面：https://域名:8162/build 根据单个 ss/ssr 信息生成calsh配置

## docker 一键集成caddy+ss+subscribe2clash方案
1、进入docker 目录

2、将.env.default 重命名为 .env 并填写域名和ss密码及加密方式，域名不填则不支持https

3、docker-compose -f docker-compose.yml up -d

4、访问https://域名:8162/build 生成连接

<img src="QQ20230708-111213.png" width="300px">
<br>
<img src="QQ20230708-111126.png" width="300px">

## 启动服务

### 二进制

- [release](https://github.com/icpd/subscribe2clash/releases)下载对应的版本
- 解压后执行`./subscribe2clash`
- 访问http://localhost:8162/?link=你的订阅链接

### 源码编译

- 安装Go 1.20+
- `go get github.com/icpd/subscribe2clash`
- 进入项目根目录
- 执行 `make build`

## 使用

### 命令行
通过命令生成 Clash 配置，方便在脚本中使用。
```bash
# 通过文件
./subscribe2clash --file=subscribe_file -o=output_file 

# 通过订阅链接
./subscribe2clash --link=subscribe_link -o=output_file

# 只输出节点信息
./subscribe2clash --link=subscribe_link -o=output_file  --nodeonly
```


### web服务
启动一个 HTTP 服务，访问 http://localhost:8162/clash?link=你的订阅链接 即可获取 Clash 配置。

访问 http://localhost:8162/build 生成单个 Clash 配置连接

```bash
./subscribe2clash
```

## 其它
- 指定自定义基础配置文件，可在里面添加自定义的路由规则，程序将按照这个文件写入路由信息。可参考[internal/acl/config/default_base_config.yaml](https://github.com/icpd/subscribe2clash/blob/master/internal/acl/config/default_base_config.yaml) ，`{{.}}`将被程序替换为ACL的路由规则。

  ```
  ./subscribe2clash -b ./yourfile.yaml
  ```

- 指定自定义碎片路由配置文件，可在该配置文件中自定义添加ACL碎片规则文件。可参考[internal/acl/config/default_rules.ini](https://github.com/icpd/subscribe2clash/blob/master/internal/acl/config/default_rules.ini)

  ```
  ./subscribe2clash -r ./yourfile.ini
  ```

- 指定输出的配置文件。默认情况下配置文件会输出为`./config/acl.yaml`，可以通过以下命令来重新指定。

  ```
  ./subscribe2clash -o ./yourconfig.yaml
  ```


- 启用http代理。由于网络原因，ACL的github源可能连接不上，你可能需要配合代理一起食用。

  ```
  ./subscribe2clash -proxy http://127.0.0.1:7890
  ```

- 指定服务监听地址，默认监听`127.0.0.1:8162`端口。

  ```
  ./subscribe2clash -l 127.0.0.1:8162
  ```

- 指定更新规则频率，单位小时，默认每6小时拉取一次。

  ```
  ./subscribe2clash -t 6
  ```