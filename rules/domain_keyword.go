/*
@Time : 2021/3/9 4:06 PM
@Author : charley
@File : domain_keyword
*/
package rules

import (
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/miekg/dns"
	"strings"
)

type DomainKeyword struct {
	keyword      string
	remoteServer constant.RemoteDnsServer
}

func (dk *DomainKeyword) RuleType() constant.RuleType {
	return constant.DomainKeyword
}

func (dk *DomainKeyword) Match(msg *dns.Msg) bool {
	if len(msg.Question) <= 0 {
		return false
	}
	domain := msg.Question[0].Name

	return strings.Contains(domain, dk.keyword)
}

func (dk *DomainKeyword) RemoteServer() constant.RemoteDnsServer {
	return dk.remoteServer
}

func (dk *DomainKeyword) ResultIPAddress() string {
	return ""
}

func (dk *DomainKeyword) Payload() string {
	return dk.keyword
}

func NewDomainKeyword(keyword string, remoteServer constant.RemoteDnsServer) *DomainKeyword {
	return &DomainKeyword{
		keyword:      keyword,
		remoteServer: remoteServer,
	}
}
