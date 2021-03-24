/*
@Time : 2021/3/12 4:00 PM
@Author : charley
@File : rule_list_config
*/
package config

type RuleListConfig struct {
	Name  string   `yaml:"name"`
	Rules []string `yaml:"rules"`
}
