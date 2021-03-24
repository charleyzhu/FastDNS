/*
@Time : 2021/3/8 5:26 PM
@Author : charley
@File : initial.go
*/
package config

import (
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/charleyzhu/FastDNS/log"
	"github.com/charleyzhu/FastDNS/utils"
	"os"
)

const (
	DefaultConfig = `

debug: true

listen:
  - port: 53
    type: udp
    minTTL: 3600
    maxTTL: 86400
	cache: 5000
    rules:
      - default-group

  - port: 53
    type: tcp
    rules:
     - default-group

forward:
  - 9.9.9.10
  - 149.112.112.10

servers:
  - name: 114DNS
    address: 114.114.114.114
    port: 53
    type: udp

  - name: aliDNS
    address: 223.5.5.5
    port: 53
    type: tcp

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
    address: 127.0.0.1
    port: 53530
    type: udp

server-group:
  - name: forward
    type: parallel
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
  - name: china
    files:
      - url: https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/accelerated-domains.china.conf
        type: dnsmasq
        server: isp

rules-list:
  - name: default
    rules:
      - DOMAIN-SUFFIX,baidu.com,isp
`
)

func Init(dir string) {
	// initial homedir
	isConfigExists, _ := utils.PathExists(dir)
	if !isConfigExists {
		if err := os.MkdirAll(dir, 0777); err != nil {
			log.Logger.Fatalf("Can't create config directory %s: %s", dir, err.Error())
		}
	}

	subscribeDir := constant.Path.SubscribeDir()
	isSubscribeExists, _ := utils.PathExists(subscribeDir)
	if !isSubscribeExists {
		if err := os.MkdirAll(subscribeDir, 0777); err != nil {
			log.Logger.Fatalf("Can't create subscribe directory %s: %s", subscribeDir, err.Error())
		}
	}

	if _, err := os.Stat(constant.Path.Config()); os.IsNotExist(err) {
		log.Logger.Info("Can't find config, create a initial config file")
		f, err := os.OpenFile(constant.Path.Config(), os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Logger.Fatalf("Can't create file %s: %s", constant.Path.Config(), err.Error())
		}
		f.Write([]byte(DefaultConfig))
		f.Close()
	}

}
