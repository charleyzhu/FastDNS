/*
@Time : 2021/3/12 3:35 PM
@Author : charley
@File : server_config
*/
package config

type ServerConfig struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	Type    string `yaml:"type"`
}
