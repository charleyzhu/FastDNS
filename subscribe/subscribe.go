/*
@Time : 2021/3/9 11:02 PM
@Author : charley
@File : subscribe
*/
package subscribe

import (
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/charleyzhu/FastDNS/log"
	"github.com/charleyzhu/FastDNS/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func Subscribe(url string) (string, error) {
	log.Logger.Infof("start subscribe: %s", url)
	cacheFilePath := GetUrlCachePath(url)
	isExist, _ := utils.PathExists(cacheFilePath)
	if isExist {
		return loadCacheFile(cacheFilePath)
	} else {
		ruleString, err := download(url)
		if err != nil {
			return "", err
		}
		return ruleString, nil
	}
}

func loadCacheFile(cacheFilePath string) (string, error) {
	log.Logger.Infof("load subscribe form cache file : %s", cacheFilePath)
	ruleString, err := ioutil.ReadFile(cacheFilePath)
	if err != nil {
		return "", err
	}
	return string(ruleString), nil
}

func SaveCacheFile(cacheFilePath string, cachePayload string) {
	f, err := os.OpenFile(cacheFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Logger.Errorf("Can't create file %s: %s", cacheFilePath, err.Error())
	}
	defer f.Close()
	f.Write([]byte(cachePayload))
}

func Update(url string) {
	cacheFilePath := GetUrlCachePath(url)
	ruleString, err := download(url)
	if err != nil {
		return
	}
	SaveCacheFile(cacheFilePath, ruleString)
}

func download(url string) (string, error) {
	log.Logger.Infof("load subscribe form url : %s", url)
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resultString := string(body)
	return resultString, nil
}

func GetUrlCachePath(url string) string {
	subscribeDir := constant.Path.SubscribeDir()
	cacheFileName := path.Base(url)
	cacheFilePath := path.Join(subscribeDir, cacheFileName)
	return cacheFilePath
}
