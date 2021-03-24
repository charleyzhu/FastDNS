/*
@Time : 2021/3/12 4:12 PM
@Author : charley
@File : listen_config
*/
package config

type ListenConfig struct {
	Port       int      `yaml:"port"`
	Type       string   `yaml:"type"`
	MinTTL     uint32   `yaml:"minTTL"`
	MaxTTL     uint32   `yaml:"maxTTL"`
	CacheCount int      `yaml:"cache"`
	Rules      []string `yaml:"rules"`
}
