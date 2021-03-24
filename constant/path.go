/*
@Time : 2021/3/8 4:54 PM
@Author : charley
@File : path
*/
package constant

import (
	"os"
	P "path"
	"path/filepath"
)

// Path is used to get the configuration path
var Path *path

type path struct {
	homeDir      string
	configFile   string
	subscribeDir string
}

func init() {

	currentDir, _ := os.Getwd()

	homeDir := P.Join(currentDir, "config")
	subscribeDir := P.Join(homeDir, "subscribe")
	Path = &path{homeDir: homeDir, configFile: "config.yaml", subscribeDir: subscribeDir}
}

// SetHomeDir is used to set the configuration path
func SetHomeDir(root string) {
	Path.homeDir = root
}

// SetConfig is used to set the configuration file
func SetConfig(file string) {
	Path.configFile = file
}

func (p *path) HomeDir() string {
	return p.homeDir
}

func (p *path) Config() string {
	return p.configFile
}

func (p *path) SubscribeDir() string {
	return p.subscribeDir
}

// Resolve return a absolute path or a relative path with homedir
func (p *path) Resolve(path string) string {
	if !filepath.IsAbs(path) {
		return filepath.Join(p.HomeDir(), path)
	}

	return path
}
