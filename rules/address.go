/*
@Time : 2021/3/9 4:07 PM
@Author : charley
@File : Address
*/
package rules

import (
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/miekg/dns"
)

type Address struct {
	domain       string
	remoteServer constant.RemoteDnsServer
}

func (addr *Address) RuleType() constant.RuleType {
	return constant.Address
}

func (addr *Address) Match(msg *dns.Msg) bool {
	if len(msg.Question) <= 0 {
		return false
	}
	domain := msg.Question[0].Name
	return addr.domain == domain
}

func (addr *Address) RemoteServer() constant.RemoteDnsServer {
	return addr.remoteServer
}

func (addr *Address) Payload() string {
	return addr.domain
}

func NewAddress(domain string, remoteServer constant.RemoteDnsServer) *Address {
	return &Address{
		domain:       domain + ".",
		remoteServer: remoteServer,
	}
}
