/*
@Time : 2021/3/12 3:22 PM
@Author : charley
@File : subscribe_file
*/
package config

type SubscribeConfig struct {
	Name  string          `yaml:"name"`
	Files []SubscribeFile `yaml:"files"`
}

type SubscribeFile struct {
	Url    string `yaml:"url"`
	Type   string `yaml:"type"`
	Server string `yaml:"server"`
}
