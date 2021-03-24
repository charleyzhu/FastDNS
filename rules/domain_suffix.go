/*
@Time : 2021/3/9 3:57 PM
@Author : charley
@File : domain_suffix
*/
package rules

import (
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/miekg/dns"
	"strings"
)

type DomainSuffix struct {
	suffix       string
	remoteServer constant.RemoteDnsServer
}

func (ds *DomainSuffix) RuleType() constant.RuleType {
	return constant.DomainSuffix
}

func (ds *DomainSuffix) Match(msg *dns.Msg) bool {
	if len(msg.Question) <= 0 {
		return false
	}
	domain := msg.Question[0].Name

	return strings.HasSuffix(domain, "."+ds.suffix) || domain == ds.suffix
}

func (ds *DomainSuffix) RemoteServer() constant.RemoteDnsServer {
	return ds.remoteServer
}

func (ds *DomainSuffix) Payload() string {
	return ds.suffix
}

func NewDomainSuffix(suffix string, remoteServer constant.RemoteDnsServer) *DomainSuffix {
	return &DomainSuffix{
		suffix:       suffix + ".",
		remoteServer: remoteServer,
	}
}
