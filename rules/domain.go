/*
@Time : 2021/3/9 4:05 PM
@Author : charley
@File : domain
*/
package rules

import (
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/miekg/dns"
)

type Domain struct {
	domain       string
	remoteServer constant.RemoteDnsServer
}

func (d *Domain) RuleType() constant.RuleType {
	return constant.Domain
}

func (d *Domain) Match(msg *dns.Msg) bool {
	if len(msg.Question) <= 0 {
		return false
	}
	domain := msg.Question[0].Name
	return d.domain == domain
}

func (d *Domain) RemoteServer() constant.RemoteDnsServer {
	return d.remoteServer
}

func (d *Domain) Payload() string {
	return d.domain
}

func NewDomain(domain string, remoteServer constant.RemoteDnsServer) *Domain {
	return &Domain{
		domain:       domain + ".",
		remoteServer: remoteServer,
	}
}
