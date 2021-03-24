/*
@Time : 2021/3/9 4:51 PM
@Author : charley
@File : rules_group
*/
package rules

import (
	"fmt"
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/miekg/dns"
)

type RuleList struct {
	name  string
	rules []constant.Rule
}

func (rl *RuleList) Rules() []constant.Rule {
	return rl.rules
}

func (rl *RuleList) AddressRules() []constant.Rule {
	var resultRules []constant.Rule
	for _, rule := range rl.rules {
		if rule.RuleType() == constant.Address ||
			rule.RuleType() == constant.AddressSuffix ||
			rule.RuleType() == constant.AddressKeyword {
			resultRules = append(resultRules, rule)
		}

	}
	return resultRules
}

func (rl *RuleList) ServerRules() []constant.Rule {
	var resultRules []constant.Rule
	for _, rule := range rl.rules {
		if rule.RuleType() == constant.Domain ||
			rule.RuleType() == constant.DomainSuffix ||
			rule.RuleType() == constant.DomainKeyword {
			resultRules = append(resultRules, rule)
		}

	}
	return resultRules
}

func (rl *RuleList) Match(msg *dns.Msg) (constant.Rule, error) {

	//先匹配地址规则
	addressRules := rl.AddressRules()
	for _, rule := range addressRules {
		if rule.Match(msg) {
			return rule, nil
		}
	}

	serverRules := rl.ServerRules()

	for _, rule := range serverRules {
		if rule.Match(msg) {
			return rule, nil
		}
	}
	return nil, fmt.Errorf("there are no matching rules")
}

func (rl *RuleList) Name() string {
	return rl.name
}

func NewRuleList(name string, rules []constant.Rule) *RuleList {
	return &RuleList{
		name:  name,
		rules: rules,
	}

}
