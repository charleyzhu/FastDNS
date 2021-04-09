/*
@Time : 2021/3/9 5:30 PM
@Author : charley
@File : parser
*/
package rules

import (
	"fmt"
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/charleyzhu/FastDNS/server/remote_server"
	"strings"
)

func ParseRule(ruleLine string, remoteDnsServers map[string]constant.RemoteDnsServer) (constant.Rule, error) {
	var (
		parseErr error
		parsed   constant.Rule
	)

	ruleArray := strings.Split(ruleLine, ",")
	if len(ruleArray) != 3 {
		return nil, fmt.Errorf("ParseRule: %s error", ruleLine)
	}
	ruleType := ruleArray[0]
	payload := ruleArray[1]
	params := strings.ToLower(ruleArray[2])

	switch ruleType {
	case "DOMAIN-SUFFIX":
		remoteServerName := params
		remoteServer, existServer := remoteDnsServers[remoteServerName]
		if existServer {
			parsed = NewDomainSuffix(payload, remoteServer)
		} else {
			parseErr = fmt.Errorf("rule Line %s server name undefined", ruleLine)
		}

	case "DOMAIN-KEYWORD":
		remoteServerName := params
		remoteServer, existServer := remoteDnsServers[remoteServerName]
		if existServer {
			parsed = NewDomainKeyword(payload, remoteServer)
		} else {
			parseErr = fmt.Errorf("rule Line %s server name undefined", ruleLine)
		}

	case "DOMAIN":
		remoteServerName := params
		remoteServer, existServer := remoteDnsServers[remoteServerName]
		if existServer {
			parsed = NewDomain(payload, remoteServer)
		} else {
			parseErr = fmt.Errorf("rule Line %s server name undefined", ruleLine)
		}

	case "ADDRESS-SUFFIX":
		remoteServer := remote_server.NewAddressDnsServer(params)
		parsed = NewAddressSuffix(payload, remoteServer)
	case "ADDRESS-KEYWORD":
		remoteServer := remote_server.NewAddressDnsServer(params)
		parsed = NewAddressKeyword(payload, remoteServer)
	case "ADDRESS":
		remoteServer := remote_server.NewAddressDnsServer(params)
		parsed = NewAddress(payload, remoteServer)
	default:
		parseErr = fmt.Errorf("unsupported rule type %s", ruleType)
	}
	return parsed, parseErr
}

func ParseClashRule(ruleLine string, remoteDnsServer constant.RemoteDnsServer) (constant.Rule, error) {
	var (
		parseErr error
		parsed   constant.Rule
	)

	ruleArray := strings.Split(ruleLine, ",")
	if len(ruleArray) != 2 {
		return nil, fmt.Errorf("ParseRule: %s error", ruleLine)
	}
	ruleType := ruleArray[0]
	payload := ruleArray[1]

	switch ruleType {
	case "DOMAIN-SUFFIX":
		parsed = NewDomainSuffix(payload, remoteDnsServer)

	case "DOMAIN-KEYWORD":
		parsed = NewDomainKeyword(payload, remoteDnsServer)

	case "DOMAIN":
		parsed = NewDomain(payload, remoteDnsServer)

	case "ADDRESS-SUFFIX":
		parsed = NewAddressSuffix(payload, remoteDnsServer)
	case "ADDRESS-KEYWORD":
		parsed = NewAddressKeyword(payload, remoteDnsServer)
	case "ADDRESS":
		parsed = NewAddress(payload, remoteDnsServer)
	default:
		parseErr = fmt.Errorf("unsupported rule type %s", ruleType)
	}
	return parsed, parseErr
}

func ParseDnsmasqRule(ruleLine string, remoteDnsServer constant.RemoteDnsServer) (constant.Rule, error) {
	ruleLineArray := strings.Split(ruleLine, "/")
	if len(ruleLineArray) != 3 {
		return nil, fmt.Errorf("rule line [%s] error", ruleLine)
	}
	ruleType := ruleLineArray[0]
	domain := ruleLineArray[1]
	params := ruleLineArray[1]

	switch strings.ToLower(ruleType) {
	case "server=":
		return NewDomainSuffix(domain, remoteDnsServer), nil
	case "address=":
		return NewAddressSuffix(domain, remote_server.NewAddressDnsServer(params)), nil
	default:
		return nil, fmt.Errorf("unknown type [%s]", ruleType)
	}

}
