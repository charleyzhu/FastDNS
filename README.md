# FastDNS


## Todo

- [ ] Auto Update Subscribe
- [ ] NOT Remote Server Support

Config file demo

```yaml
# 是否开启调试模式
debug: true

listen:
  - port: 53  # 监听端口
    type: udp # 监听协议
    minTTL: 172800 #最小ttl（可选） 如果上级dns小于这个值那么修改为这个值
    maxTTL: 259200 #最大ttl（可选） 使用规则如上
    cache: 5000 # 缓存条数（可选）
    rules: # 规则数组 可以是rules-group，rules-subscribe，rules-list中的name组合
      - default-group

  - port: 53
    type: tcp
    rules:
     - default-group

# 上级dns 列表只能是udp服务器
forward: 
  - 9.9.9.10
  - 149.112.112.10

#服务器列表
servers:
  - name: 114DNS  # 名称
    address: 114.114.114.114 # 服务器地址
    port: 53 # 服务器端口
    type: udp # 服务器协议

  - name: aliDNS
    address: 223.5.5.5
    port: 53
    type: tcp

  - name: tencent
    address: dot.pub
    port: 443
    type: tls

  - name: google1
    address: 8.8.8.8
    port: 53
    type: tcp

  - name: google2
    address: 8.8.4.4
    port: 53
    type: tcp

  - name: isp1
    address: 59.51.78.211
    port: 53
    type: udp

  - name: isp2
    address: 222.246.129.81
    port: 53
    type: udp

  - name: clash
    address: 172.16.100.2
    port: 53
    type: udp

# 服务器组
server-group:
  - name: forward #服务器组 名称 如果命名为forward就会覆盖默认的forward数组，可以定义他们的负载模式
    type: parallel #负载模式 parallel：一个一个服务器的查询，直到查询到结果，balancing，异步向所有服务器发送请求，等待最先返回的数据，fasttest检测所有服务器返回的a，记录返回ping最小的a记录
    servers:
      - clash

  - name: isp
    type: parallel
    servers:
      - isp1
      - isp2

  - name: public-dns
    type: balancing
    servers:
      - 114DNS
      - aliDNS

  - name: Balancing
    type: balancing
    servers:
      - 114DNS
      - aliDNS

  - name: Parallel
    type: parallel
    servers:
      - 114DNS
      - aliDNS

  - name: FastTest
    type: fasttest
    servers:
      - 114DNS
      - aliDNS

rules-group:
  - name: default-group
    rules:
      - default
      - china

rules-subscribe:
  # 订阅名称
  - name: china
    # 订阅文件
    files:
      # 订阅的地址
      # - url: https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/accelerated-domains.china.conf
      # 订阅类型 如果订阅内容和rules-list一样则不需要填写type字断
      #   type: dnsmasq
      # 订阅适配服务器（服务器组） 如果订阅内容和rules-list一样则不需要填写server字断
      #   server: isp

rules-list:
  # 规则列表
  - name: default
    rules:
      # 规则 共有6种定义
      # DOMAIN 完整适配域名
      # DOMAIN-SUFFIX 适配域名后缀
      # DOMAIN-KEYWORD 适配域名关键字
      # 例子- DOMAIN-SUFFIX,0-100.com,isp 
      # 规则方式, 规则表达式,服务器名称或服务器组名称
      
      # ADDRESS 完整适配域名 返回指定地址
      # ADDRESS-SUFFIX 适配域名后缀 返回指定地址
      # ADDRESS-KEYWORD 适配域名关键字 返回指定地址
      # 例子- DOMAIN-SUFFIX,0-100.com,192.168.1.100 
      # 规则方式, 规则表达式,服务器名称或服务器组名称
      - DOMAIN,www.baidu.com,isp
      - DOMAIN-SUFFIX,baidu.com,isp
      - DOMAIN-KEYWORD,baidu,isp
      - ADDRESS,www.baidu.com,192.168.1.100
      - ADDRESS-SUFFIX,baidu.com,isp
      - ADDRESS-KEYWORD,baidu,isp
      
```