/*
@Time : 2021/3/8 4:41 PM
@Author : charley
@File : rule
*/
package constant

import (
	"github.com/miekg/dns"
)

const (
	Domain RuleType = iota
	DomainSuffix
	DomainKeyword
	Address
	AddressSuffix
	AddressKeyword
)

type RuleType int

func (rt RuleType) String() string {
	switch rt {
	case Domain:
		return "Domain"
	case DomainSuffix:
		return "DomainSuffix"
	case DomainKeyword:
		return "DomainKeyword"
	case Address:
		return "Address"
	case AddressSuffix:
		return "AddressSuffix"
	case AddressKeyword:
		return "AddressKeyword"
	default:
		return "Unknown"
	}
}

type Rule interface {
	RuleType() RuleType
	Match(msg *dns.Msg) bool
	RemoteServer() RemoteDnsServer
	Payload() string
}

type RulesGroup interface {
	Rules() []Rule
	Name() string
	Match(msg *dns.Msg) (Rule, error)
}
