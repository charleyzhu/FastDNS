/*
@Time : 2021/3/12 4:05 PM
@Author : charley
@File : rule_group_config
*/
package config

type RuleGroupConfig struct {
	Name  string   `yaml:"name"`
	Rules []string `yaml:"rules"`
}
