agile-proxy
======
一款灵活，轻量，易于扩展的代理工具。

### 下载
***
https://github.com/fanyiguang/agile-proxy/releases/tag/v1.0.0

### 功能
***
* 支持HTTP/HTTPS/SOCKS5/SSL/SSH代理协议
* 多种协议端口监听
* 多协议代理链
* 模拟动态代理
* 自定义配置代理链路
* 高可用代理模式
* 自定义DNS
* 消息处理

### 快速上手
***
##### 逻辑结构
![image](https://github.com/fanyiguang/agile-proxy/blob/main/_example/res/agile_proxy.png)

##### 自定义链路 
![image](https://github.com/fanyiguang/agile-proxy/blob/main/_example/res/muti_link.png)  
每个模块下都有多种实现。它们像组件一般，您可以挑选合适的组件来拼装出属于自己的代理链路。

##### 配置SOCKS5服务端
```bash
{
  "ip": "127.0.0.1",
  "port": "7890",
  "username": "admin",
  "password": "123",
  "type": "socks5",
  "name": "socks5-server",
  "auth_mode": 1, //认证模式 0-只允许认证模式 1-允许匿名模式
  "transport_name": "direct-transport" //传输器模块名称
}
```

##### 配置DIRECT传输器
```bash
{
  "type": "direct",
  "name": "direct-transport",
  "client_name": "socks5-client",
  "dns": {
    "local_dns": false,
    "server": "114.114.114.114"
  }
}
```

##### 配置SOCKS5客户端
```bash
{
  "ip": "127.0.0.1",
  "port": "7892",
  "username": "admin",
  "password": "123",
  "type": "socks5",
  "name": "socks5-client",
  "auth_mode": 1,
  "dialer_name": "direct-dialer", //连接器模块名称
  "mode": 0 // 连接模式 0-降级模式 1-严格模式
}
```
在降级模式下如果连接器拨号失败了，客户端会切换到默认连接模式再次尝试连接目标服务器，严格模式的话则不会尝试默认连接。

##### 配置DIRECT连接器
```bash
{
  "type": "direct",
  "name": "direct-dialer",
  "interface": "" //本地网卡ip
}
```

##### 更多配置说明
https://github.com/fanyiguang/agile-proxy/tree/main/_example/config_description

### 高级应用
***
##### 配置动态传输器
```bash
{
  "type": "dynamic",
  "name": "dynamic",
  "client_names": "socks5-1,socks5-2,direct", //多个客户端用,隔开
  "rand_rule": "timestamp", //时间戳(默认)
  "dns": {
    "local_dns": false,
    "server": "114.114.114.114"
  }
}
```

##### 配置高可用传输器
```bash
{
  "type": "ha",
  "name": "ha",
  "client_names": "socks5-1,socks5-2",
  "dns": {
    "local_dns": false,
    "server": "114.114.114.114"
  }
}
```

### 更多配置示例
***
https://github.com/fanyiguang/agile-proxy/tree/main/_example/config

### 定制自己的代理
***
个人的力量是微不足道的，群众的力量才是无限的。agile-proxy设计之初的理念就是高扩展性，项目中几乎所有功能都是完全解耦的。您可以很容易的接入开发创造出属于自己的代理。





