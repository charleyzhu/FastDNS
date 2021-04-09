/*
@Time : 2021/3/8 4:39 PM
@Author : charley
@File : config
*/
package config

import (
	"fmt"
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/charleyzhu/FastDNS/log"
	"github.com/charleyzhu/FastDNS/rules"
	"github.com/charleyzhu/FastDNS/server/local_server"
	"github.com/charleyzhu/FastDNS/server/remote_server"
	"github.com/charleyzhu/FastDNS/subscribe"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	Listen           []local_server.LocalServer
	RemoteDnsServers map[string]constant.RemoteDnsServer
	RulesList        map[string]constant.RulesGroup
	RuleSubscribe    map[string]constant.RulesGroup
	RuleGroup        map[string]constant.RulesGroup
}

type RawConfig struct {
	Debug         bool                `yaml:"debug"`
	Listen        []ListenConfig      `yaml:"listen"`
	ForwardServer []string            `yaml:"forward"`
	Servers       []ServerConfig      `yaml:"servers"`
	ServerGroup   []ServerGroupConfig `yaml:"server-group"`

	RuleGroup     []RuleGroupConfig `yaml:"rules-group"`
	RuleSubscribe []SubscribeConfig `yaml:"rules-subscribe"`
	RuleList      []RuleListConfig  `yaml:"rules-list"`
}

// 获得配置
func GetConfig() (*Config, error) {
	rawCfg, err := ParseWithPath(constant.Path.Config())

	if err != nil {
		return nil, err
	}
	log.InitLogger(rawCfg.Debug)
	return ParseRawConfig(rawCfg)
}

//解析配置文件
func ParseWithPath(path string) (*RawConfig, error) {
	buf, err := readConfig(path)
	if err != nil {
		return nil, err
	}
	config := &RawConfig{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, err
	}
	return config, nil

}

//读取配置文件
func readConfig(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("configuration file %s is empty", path)
	}

	return data, err
}

func ParseRawConfig(rawCfg *RawConfig) (*Config, error) {
	config := &Config{}

	remoteDnsServers, err := parseRemoteServers(rawCfg)
	if err != nil {
		return nil, err
	}
	_, existForward := remoteDnsServers["forward"]
	if !existForward {
		forwardServerGroup, err := parseForwardDns(rawCfg)
		if err != nil {
			return nil, err
		}
		remoteDnsServers["forward"] = forwardServerGroup
	}
	config.RemoteDnsServers = remoteDnsServers

	rulesList, err := parseRules(rawCfg, remoteDnsServers)
	if err != nil {
		return nil, err
	}
	config.RulesList = rulesList

	ruleSubscribe, err := parseSubscribe(rawCfg, remoteDnsServers)
	if err != nil {
		config.RuleSubscribe = make(map[string]constant.RulesGroup)
		log.Logger.Error(err)
	} else {
		// 检测重复名称
		err := duplicateNames(rulesList, ruleSubscribe)
		if err != nil {
			return nil, err
		}
		config.RuleSubscribe = ruleSubscribe
	}

	rulesGroup, err := parseRulesGroup(rawCfg, config.RulesList, config.RuleSubscribe)
	if err != nil {
		return nil, err
	}
	// 检测重复名称
	err = duplicateNames(rulesList, rulesGroup)
	if err != nil {
		return nil, err
	}
	// 检测重复名称
	err = duplicateNames(ruleSubscribe, rulesGroup)
	if err != nil {
		return nil, err
	}

	config.RuleGroup = rulesGroup

	listen, err := parseListen(rawCfg, *config)
	if err != nil {
		return nil, err
	}
	config.Listen = listen

	return config, nil
}

// 检测重复名称
func duplicateNames(mapA map[string]constant.RulesGroup, mapB map[string]constant.RulesGroup) error {
	for aKey := range mapA {
		for bKey := range mapB {
			if aKey == bKey {
				return fmt.Errorf("duplicate rule name:%s", aKey)
			}
		}
	}
	return nil
}

// 解析默认dns
func parseForwardDns(rawCfg *RawConfig) (*remote_server.DnsServerGroup, error) {

	forwardDnsCfg := rawCfg.ForwardServer
	var groupServers []constant.RemoteDnsServer
	for _, line := range forwardDnsCfg {
		server := remote_server.NewUDPDnsServer(line, line, 53)
		groupServers = append(groupServers, server)
	}

	forwardServers := remote_server.NewDnsServerGroup("forward", remote_server.Balancing, groupServers)
	return forwardServers, nil
}

// 解析远程dns服务器
func parseRemoteServers(rawCfg *RawConfig) (map[string]constant.RemoteDnsServer, error) {
	servers := rawCfg.Servers

	remoteServersMap := make(map[string]constant.RemoteDnsServer)

	for idx, serverConfig := range servers {

		if serverConfig.Name == "" {
			return nil, fmt.Errorf("servers %d: name error", idx)
		}

		lowerName := strings.ToLower(serverConfig.Name)

		if serverConfig.Address == "" {
			return nil, fmt.Errorf("servers %d: Address error", idx)
		}

		if serverConfig.Port == 0 {
			return nil, fmt.Errorf("servers %d: port error", idx)
		}

		if serverConfig.Type == "" {
			return nil, fmt.Errorf("servers %d: type error", idx)
		}

		//检查是否重名
		_, existName := remoteServersMap[lowerName]
		if existName {
			return nil, fmt.Errorf("servers %d: Duplicate name:%s", idx, serverConfig.Name)
		}

		remoteServer, err := remote_server.ParseRemoteServer(serverConfig.Type, serverConfig.Name, serverConfig.Address, serverConfig.Port)
		if err != nil {
			return nil, err
		}

		remoteServersMap[lowerName] = remoteServer
	}

	serverGROUP := rawCfg.ServerGroup

	for idx, groupConfig := range serverGROUP {

		if groupConfig.Name == "" {
			return nil, fmt.Errorf("serverConfig-groupConfig %d: name error", idx)
		}

		lowerName := strings.ToLower(groupConfig.Name)

		if groupConfig.GroupType == "" {
			return nil, fmt.Errorf("serverConfig-groupConfig %d: type error", idx)
		}

		if groupConfig.Servers == nil {
			return nil, fmt.Errorf("serverConfig-groupConfig %d: servers error", idx)
		}

		//检查是否重名
		_, existName := remoteServersMap[lowerName]
		if existName {
			return nil, fmt.Errorf("servers-groupConfig %d: Duplicate name:%s", idx, groupConfig.Name)
		}

		var gType remote_server.DnsServerGroupType
		switch groupConfig.GroupType {
		case "balancing":
			gType = remote_server.Balancing
		case "parallel":
			gType = remote_server.Parallel
		case "fasttest":
			gType = remote_server.FastTest
		default:
			return nil, fmt.Errorf("serverConfig-groupConfig %d: type error", idx)
		}

		var groupServers []constant.RemoteDnsServer
		for _, serverName := range groupConfig.Servers {
			serverName = strings.ToLower(serverName)
			dnsServer, existServerName := remoteServersMap[serverName]
			if !existServerName {
				return nil, fmt.Errorf("serverConfig-groupConfig %d: serverName:%s exist servers", idx, serverName)
			}
			groupServers = append(groupServers, dnsServer)
		}

		serverGroup := remote_server.NewDnsServerGroup(groupConfig.Name, gType, groupServers)

		remoteServersMap[lowerName] = serverGroup
	}

	return remoteServersMap, nil
}

//解析规则列表
func parseRules(rawCfg *RawConfig, remoteDnsServers map[string]constant.RemoteDnsServer) (map[string]constant.RulesGroup, error) {
	ruleList := rawCfg.RuleList
	rulesGroup := make(map[string]constant.RulesGroup)

	for listIdx, ruleConfig := range ruleList {

		if ruleConfig.Name == "" {
			return nil, fmt.Errorf("rules-list %d: name error", listIdx)
		}

		lowerName := strings.ToLower(ruleConfig.Name)

		if ruleConfig.Rules == nil {
			return nil, fmt.Errorf("rules-list %d: servers error", listIdx)
		}

		//检查是否重名
		_, existName := rulesGroup[lowerName]
		if existName {
			return nil, fmt.Errorf("rules-list %d: Duplicate name:%s", listIdx, ruleConfig.Name)
		}

		var rulesListArray []constant.Rule

		for _, ruleLine := range ruleConfig.Rules {

			rule, err := rules.ParseRule(ruleLine, remoteDnsServers)
			if err != nil {
				return nil, err
			}
			rulesListArray = append(rulesListArray, rule)
		}
		rulesGroup[lowerName] = rules.NewRuleList(ruleConfig.Name, rulesListArray)
	}

	return rulesGroup, nil
}

// 解析订阅规则列表
func parseSubscribe(rawCfg *RawConfig, remoteDnsServers map[string]constant.RemoteDnsServer) (map[string]constant.RulesGroup, error) {
	ruleListCfg := rawCfg.RuleSubscribe

	rulesGroupMap := make(map[string]constant.RulesGroup)

	rulesGroupWaitGroup := new(sync.WaitGroup)
	chanRulesGroup := make(chan constant.RulesGroup, len(ruleListCfg))
	rulesGroupWaitGroup.Add(len(ruleListCfg))

	for subscribeListIdx, ruleSubscribeConfig := range ruleListCfg {

		if ruleSubscribeConfig.Name == "" {
			return nil, fmt.Errorf("rules-subscribe %d: name error", subscribeListIdx)
		}

		//if ruleSubscribeConfig.Files == nil {
		//	log.Logger.Warnf("rules-subscribe %d: No definition files ", subscribeListIdx)
		//	rulesGroupWaitGroup.Done()
		//	continue
		//} else {
		go parseSubscribeFile(chanRulesGroup, ruleSubscribeConfig, remoteDnsServers, rulesGroupWaitGroup)
		//}

	}
	rulesGroupWaitGroup.Wait()
	close(chanRulesGroup)

	for ruleGroup := range chanRulesGroup {

		ruleGroupName := ruleGroup.Name()
		//检查是否重名
		_, existName := rulesGroupMap[ruleGroupName]
		if existName {
			return nil, fmt.Errorf("rules-subscribe Duplicate name:%s", ruleGroupName)
		} else {
			rulesGroupMap[ruleGroupName] = ruleGroup
		}

	}

	return rulesGroupMap, nil
}

func parseSubscribeFile(chanRulesGroup chan constant.RulesGroup, subscribeConfig SubscribeConfig, remoteDnsServers map[string]constant.RemoteDnsServer, rulesGroupWaitGroup *sync.WaitGroup) {

	lowerName := strings.ToLower(subscribeConfig.Name)
	var ruleList []constant.Rule
	if subscribeConfig.Files == nil {
		log.Logger.Warnf("rules-subscribe %s: No definition files ", lowerName)
		chanRulesGroup <- rules.NewRuleList(lowerName, ruleList)
		rulesGroupWaitGroup.Done()
		return
	}

	filesArray := subscribeConfig.Files

	ruleChan := make(chan []constant.Rule, len(filesArray))
	ruleListWaitGroup := new(sync.WaitGroup)
	ruleListWaitGroup.Add(len(filesArray))

	for fileIdx, subscribeFile := range filesArray {

		if subscribeFile.Url == "" {
			log.Logger.Errorf("rules-subscribe [%s] [%d]: missing url", subscribeConfig.Name, fileIdx)
			ruleListWaitGroup.Done()
			break
		}

		if subscribeFile.Server == "" {
			log.Logger.Errorf("rules-subscribe [%s] [%d]: missing Server", subscribeConfig.Name, fileIdx)
			ruleListWaitGroup.Done()
			break
		}

		if subscribeFile.Type != "" {
			//file type 字断存在
			if subscribeFile.Type == "dnsmasq" {
				// dnsmasq 规则解析
				remoteServer, existRemoteServer := remoteDnsServers[subscribeFile.Server]
				if !existRemoteServer {
					log.Logger.Errorf("rules-subscribe [%s] [%d]: exist RemoteServer %s", subscribeConfig.Name, fileIdx, subscribeFile.Server)
					break
				}
				go parseDnsmasqRules(subscribeFile.Url, remoteServer, ruleChan, ruleListWaitGroup)

			} else if subscribeFile.Type == "clash" {
				remoteServer, existRemoteServer := remoteDnsServers[subscribeFile.Server]
				if !existRemoteServer {
					log.Logger.Errorf("rules-subscribe [%s] [%d]: exist RemoteServer %s", subscribeConfig.Name, fileIdx, subscribeFile.Server)
					break
				}
				go parseClashRules(subscribeFile.Url, remoteServer, ruleChan, ruleListWaitGroup)
			} else {
				go parseSubscribeRules(subscribeFile.Url, remoteDnsServers, ruleChan, ruleListWaitGroup)
			}
		} else {
			// file type 没有定义说明是FastDns规则
			go parseSubscribeRules(subscribeFile.Url, remoteDnsServers, ruleChan, ruleListWaitGroup)
		}

	}

	ruleListWaitGroup.Wait()
	close(ruleChan)
	for rs := range ruleChan {
		ruleList = append(ruleList, rs...)
	}
	rulesGroupWaitGroup.Done()
	chanRulesGroup <- rules.NewRuleList(lowerName, ruleList)

}

func parseClashRules(url string, remoteDnsServer constant.RemoteDnsServer, ruleChan chan []constant.Rule, ruleListWaitGroup *sync.WaitGroup) {
	log.Logger.Infof("start parse clash rules %s", url)
	var ruleList []constant.Rule
	rulesString, err := subscribe.Subscribe(url)
	if err != nil {
		log.Logger.Error(err)
		ruleChan <- ruleList
		ruleListWaitGroup.Done()
		return
	}

	clashRuleFile := &ClashRule{}
	err = yaml.Unmarshal([]byte(rulesString), clashRuleFile)
	if err != nil {
		log.Logger.Error(err)
		ruleChan <- ruleList
		ruleListWaitGroup.Done()
		return
	}

	for _, ruleLine := range clashRuleFile.Payload {
		rule, err := rules.ParseClashRule(ruleLine, remoteDnsServer)
		if err != nil {
			log.Logger.Warn(err)
			continue
		}
		ruleList = append(ruleList, rule)
	}
	if len(ruleList) > 0 {
		savePath := subscribe.GetUrlCachePath(url)
		subscribe.SaveCacheFile(savePath, rulesString)
	}
	ruleChan <- ruleList
	ruleListWaitGroup.Done()
}

func parseDnsmasqRules(url string, remoteDnsServer constant.RemoteDnsServer, ruleChan chan []constant.Rule, ruleListWaitGroup *sync.WaitGroup) {
	log.Logger.Infof("start parse dnsmasq rules %s", url)
	var ruleList []constant.Rule
	rulesString, err := subscribe.Subscribe(url)
	if err != nil {
		log.Logger.Error(err)
		ruleChan <- ruleList
		ruleListWaitGroup.Done()
		return
	}

	rulesStringArray := strings.Split(rulesString, "\n")
	for _, ruleLine := range rulesStringArray {
		rule, err := rules.ParseDnsmasqRule(ruleLine, remoteDnsServer)
		if err != nil {
			log.Logger.Error(err)
			continue

		}
		ruleList = append(ruleList, rule)
	}
	if len(ruleList) > 0 {
		savePath := subscribe.GetUrlCachePath(url)
		subscribe.SaveCacheFile(savePath, rulesString)
	}
	ruleChan <- ruleList
	ruleListWaitGroup.Done()
}

func parseSubscribeRules(url string, remoteDnsServers map[string]constant.RemoteDnsServer, ruleChan chan []constant.Rule, ruleListWaitGroup *sync.WaitGroup) {
	log.Logger.Infof("start parse fast dns rules %s", url)
	var ruleList []constant.Rule

	rulesString, err := subscribe.Subscribe(url)
	if err != nil {
		log.Logger.Error(err)
		ruleChan <- ruleList
		ruleListWaitGroup.Done()
		return
	}
	rulesStringArray := strings.Split(rulesString, "\n")
	for _, ruleLine := range rulesStringArray {
		rule, err := rules.ParseRule(ruleLine, remoteDnsServers)
		if err != nil {
			log.Logger.Error(err)
			break
			//return nil, err
		}
		ruleList = append(ruleList, rule)
	}
	ruleChan <- ruleList
	ruleListWaitGroup.Done()
}

// 解析规则组
func parseRulesGroup(rawCfg *RawConfig, rulesList, ruleSubscribe map[string]constant.RulesGroup) (map[string]constant.RulesGroup, error) {
	ruleGroup := rawCfg.RuleGroup
	rulesGroup := make(map[string]constant.RulesGroup)

	for groupIdx, ruleGroupConfig := range ruleGroup {

		if ruleGroupConfig.Name == "" {
			return nil, fmt.Errorf("rules-group %d: name error", groupIdx)
		}

		lowerName := strings.ToLower(ruleGroupConfig.Name)

		if ruleGroupConfig.Rules == nil {
			return nil, fmt.Errorf("rules-group %d: servers error", groupIdx)
		}

		//检查是否重名
		_, existName := rulesGroup[lowerName]
		if existName {
			return nil, fmt.Errorf("rules-group %d: Duplicate name:%s", groupIdx, ruleGroupConfig.Name)
		}

		var rulesListArray []constant.Rule
		for _, ruleNameLine := range ruleGroupConfig.Rules {

			ruleNameLine = strings.ToLower(ruleNameLine)

			rulesList, existRulesInList := rulesList[ruleNameLine]
			rulesSubscribeList, existRulesInSubscribeList := ruleSubscribe[ruleNameLine]
			if !existRulesInList && !existRulesInSubscribeList {
				return nil, fmt.Errorf("rules-group %d RuleList name: %s Rule List Exist", groupIdx, ruleNameLine)
			}
			if existRulesInList {
				rs := rulesList.Rules()
				rulesListArray = append(rulesListArray, rs...)
			}

			if existRulesInSubscribeList {
				rs := rulesSubscribeList.Rules()
				rulesListArray = append(rulesListArray, rs...)
			}

		}
		rulesGroup[lowerName] = rules.NewRuleList(ruleGroupConfig.Name, rulesListArray)
	}
	return rulesGroup, nil
}

//解析本地监听
func parseListen(rawCfg *RawConfig, config Config) ([]local_server.LocalServer, error) {
	listenCfg := rawCfg.Listen
	var listenServerList []local_server.LocalServer

	for listenIdx, listenConfig := range listenCfg {

		if listenConfig.Port == 0 {
			return nil, fmt.Errorf("listen %d: port error", listenIdx)
		}
		portString := strconv.Itoa(listenConfig.Port)

		if listenConfig.MaxTTL == 0 {
			listenConfig.MaxTTL = constant.DefaultMaxTTL
		}

		if listenConfig.MinTTL == 0 {
			listenConfig.MaxTTL = constant.DefaultMinTTL
		}

		if listenConfig.MaxTTL <= listenConfig.MinTTL {
			return nil, fmt.Errorf("listen %d MaxTTL <= MinTTL", listenIdx)
		}

		if listenConfig.Type == "" {
			return nil, fmt.Errorf("listen %d: type error", listenIdx)
		}

		typeString := strings.ToLower(listenConfig.Type)
		if typeString != "udp" && typeString != "tcp" {
			return nil, fmt.Errorf("listen %d: type Unknown", listenIdx)
		}

		if listenConfig.Rules == nil {
			return nil, fmt.Errorf("rules-group %d: rules error", listenIdx)
		}

		var rulesListArray []constant.Rule
		for _, ruleNameLine := range listenConfig.Rules {

			rulesList, existRulesInList := config.RulesList[ruleNameLine]
			rulesSubscribeList, existRulesInSubscribeList := config.RuleSubscribe[ruleNameLine]
			rulesGroupList, existRulesInGroupList := config.RuleGroup[ruleNameLine]
			if !existRulesInList && !existRulesInSubscribeList && !existRulesInGroupList {
				return nil, fmt.Errorf("listen %d RuleList name %s: Rule List Exist", listenIdx, ruleNameLine)
			}

			if existRulesInList {
				rs := rulesList.Rules()
				rulesListArray = append(rulesListArray, rs...)
			}
			if existRulesInSubscribeList {
				rs := rulesSubscribeList.Rules()
				rulesListArray = append(rulesListArray, rs...)
			}

			if existRulesInGroupList {
				rs := rulesGroupList.Rules()
				rulesListArray = append(rulesListArray, rs...)
			}

		}

		ruleList := rules.NewRuleList(portString+typeString, rulesListArray)

		var (
			parseErr error
			parsed   local_server.LocalServer
		)

		forwardServer := config.RemoteDnsServers["forward"]

		switch typeString {
		case "udp":
			parsed = local_server.NewLocalServer(portString, local_server.LocalServerUDP,
				listenConfig.MinTTL, listenConfig.MaxTTL, listenConfig.CacheCount, ruleList, forwardServer)
		case "tcp":
			parsed = local_server.NewLocalServer(portString, local_server.LocalServerTCP,
				listenConfig.MinTTL, listenConfig.MaxTTL, listenConfig.CacheCount, ruleList, forwardServer)
		default:
			parseErr = fmt.Errorf("unsupported rule type %s", typeString)
		}
		if parseErr != nil {
			return nil, parseErr
		} else {
			listenServerList = append(listenServerList, parsed)
		}
	}
	return listenServerList, nil
}
