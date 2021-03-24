/*
@Time : 2021/3/12 3:48 PM
@Author : charley
@File : server_group_config
*/
package config

type ServerGroupConfig struct {
	Name      string   `yaml:"name"`
	GroupType string   `yaml:"type"`
	Servers   []string `yaml:"servers"`
}
