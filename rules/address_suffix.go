/*
@Time : 2021/3/9 4:26 PM
@Author : charley
@File : address_suffix
*/
package rules

import (
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/miekg/dns"
	"strings"
)

type AddressSuffix struct {
	suffix       string
	remoteServer constant.RemoteDnsServer
}

func (as *AddressSuffix) RuleType() constant.RuleType {
	return constant.AddressSuffix
}

func (as *AddressSuffix) Match(msg *dns.Msg) bool {
	if len(msg.Question) <= 0 {
		return false
	}
	domain := msg.Question[0].Name

	return strings.HasSuffix(domain, "."+as.suffix) || domain == as.suffix
}

func (as *AddressSuffix) RemoteServer() constant.RemoteDnsServer {
	return as.remoteServer
}

func (as *AddressSuffix) Payload() string {
	return as.suffix
}

func NewAddressSuffix(suffix string, remoteServer constant.RemoteDnsServer) *AddressSuffix {
	return &AddressSuffix{
		suffix:       suffix + ".",
		remoteServer: remoteServer,
	}
}
