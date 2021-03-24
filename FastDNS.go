/*
@Time : 2021/3/8 11:30 AM
@Author : charley
@File : FastDNS
*/
package main

import (
	"flag"
	"fmt"
	"github.com/charleyzhu/FastDNS/config"
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/charleyzhu/FastDNS/log"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
)

var (
	homeDir    string
	configFile string
	//clearCache string
)

func init() {
	flag.StringVar(&homeDir, "d", "", "set configuration directory")
	flag.StringVar(&configFile, "f", "", "specify configuration file")
	//flag.StringVar(&clearCache, "clearCache", "all", "specify configuration file")
	flag.Parse()
}

func main() {

	for _, s := range os.Args {
		if s == "clearCache" {
			clearCache()
			return
		}
	}

	if homeDir != "" {
		if !filepath.IsAbs(homeDir) {
			currentDir, _ := os.Getwd()
			homeDir = filepath.Join(currentDir, homeDir)
		}
		constant.SetHomeDir(homeDir)
	}

	if configFile != "" {
		if !filepath.IsAbs(configFile) {
			currentDir, _ := os.Getwd()
			configFile = filepath.Join(currentDir, configFile)
		}
		constant.SetConfig(configFile)
	} else {
		configFile := filepath.Join(constant.Path.HomeDir(), constant.Path.Config())
		constant.SetConfig(configFile)
	}

	config.Init(constant.Path.HomeDir())
	conf, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Get Config Error: %s", err.Error())

		return
	}

	for _, localserver := range conf.Listen {
		go localserver.ListenAndServe()
	}

	sig := make(chan os.Signal)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	log.Logger.Fatalf("Signal (%v) received, stopping\n", s)

}

func clearCache() {
	subDir := constant.Path.SubscribeDir()
	dir, err := ioutil.ReadDir(subDir)
	if err != nil {
		log.Logger.Error(err)
	}
	for _, d := range dir {
		err := os.RemoveAll(path.Join(subDir, d.Name()))
		if err != nil {
			log.Logger.Error(err)
		}
	}
	log.Logger.Info("clear success")
}
