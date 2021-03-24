/*
@Time : 2021/3/9 4:26 PM
@Author : charley
@File : address_keyword
*/
package rules

import (
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/miekg/dns"
	"strings"
)

type AddressKeyword struct {
	keyword      string
	remoteServer constant.RemoteDnsServer
}

func (ak *AddressKeyword) RuleType() constant.RuleType {
	return constant.AddressKeyword
}

func (ak *AddressKeyword) Match(msg *dns.Msg) bool {
	if len(msg.Question) <= 0 {
		return false
	}
	domain := msg.Question[0].Name

	return strings.Contains(domain, ak.keyword)
}

func (ak *AddressKeyword) RemoteServer() constant.RemoteDnsServer {
	return ak.remoteServer
}

func (ak *AddressKeyword) Payload() string {
	return ak.keyword
}

func NewAddressKeyword(keyword string, remoteServer constant.RemoteDnsServer) *AddressKeyword {
	return &AddressKeyword{
		keyword:      keyword,
		remoteServer: remoteServer,
	}
}
