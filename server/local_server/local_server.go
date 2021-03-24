/*
@Time : 2021/3/8 4:49 PM
@Author : charley
@File : local_server.go
*/
package local_server

import (
	"fmt"
	"github.com/bluele/gcache"
	"github.com/charleyzhu/FastDNS/constant"
	"github.com/charleyzhu/FastDNS/log"
	"github.com/charleyzhu/FastDNS/rules"
	"github.com/miekg/dns"
	"sort"
	"time"
)

const (
	LocalServerTCP LocalServerType = iota
	LocalServerUDP
)

type LocalServerType int

func (lst LocalServerType) String() string {
	switch lst {
	case LocalServerTCP:
		return "tcp"
	case LocalServerUDP:
		return "udp"
	default:
		return "udp"
	}
}

type LocalServer struct {
	Port       string
	ServerType LocalServerType
	MinTTL     uint32
	MaxTTL     uint32
	RulesList  *rules.RuleList
	forward    constant.RemoteDnsServer
	cacheCount int
	cache      gcache.Cache
}

func (ls *LocalServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	cacheMsg, err := ls.GetCache(r)
	if err == nil {
		// 有缓存
		log.Logger.Debugf("match Cache :%s Question Type:%d", r.Question[0].Name, r.Question[0].Qtype)
		cacheMsg.SetReply(r)
		ls.CheckTTL(cacheMsg)
		w.WriteMsg(cacheMsg)
		return
	}
	// 没缓存继续查询匹配
	remoteServer, err := ls.MatchRemoteServer(r)
	if err != nil {
		msg, _, err := ls.forward.Query(r)
		if err != nil {
			emptyMsg := dns.Msg{}
			emptyMsg.SetReply(r)
			ls.CheckTTL(msg)
			w.WriteMsg(&emptyMsg)
			log.Logger.Error(err)
			return
		} else {
			err := ls.SaveCache(msg)
			if err != nil {
				log.Logger.Error(err)
			}
			ls.CheckTTL(msg)
			w.WriteMsg(msg)
		}

	} else {

		msg, _, err := remoteServer.Query(r)
		if err != nil {
			emptyMsg := dns.Msg{}
			emptyMsg.SetReply(r)
			ls.CheckTTL(msg)
			w.WriteMsg(&emptyMsg)
			log.Logger.Error(err)
			return
		}
		if remoteServer.DnsServerType() != constant.AddressServer {
			err = ls.SaveCache(msg)
			if err != nil {
				log.Logger.Error(err)
			}
		}
		ls.CheckTTL(msg)
		w.WriteMsg(msg)

	}

}

func (ls *LocalServer) GetCache(msg *dns.Msg) (*dns.Msg, error) {

	if len(msg.Question) <= 0 {
		return nil, fmt.Errorf("question length error")
	}
	key := msg.Question[0].String()

	value, err := ls.cache.Get(key)
	if err != nil {
		return nil, err
	}
	cacheMsg, ok := value.(*dns.Msg)
	if !ok {
		log.Logger.Errorf("cache type error")
		return nil, fmt.Errorf("cache type error")
	}
	return cacheMsg, nil
}

func (ls *LocalServer) SaveCache(msg *dns.Msg) error {

	if len(msg.Question) <= 0 {
		return fmt.Errorf("question length error")
	}
	key := msg.Question[0].String()

	answer := msg.Answer
	if len(msg.Answer) <= 0 {
		log.Logger.Debugf("answer length 0: type: %s skip cache", dns.Type(msg.Question[0].Qtype).String())
		return nil
	}

	sort.Slice(answer, func(i, j int) bool {
		return answer[i].Header().Ttl > answer[i].Header().Ttl
	})
	expiration := time.Duration(answer[0].Header().Ttl) * time.Second
	return ls.cache.SetWithExpire(key, msg, expiration)
}

func (ls *LocalServer) CheckTTL(msg *dns.Msg) {
	for _, a := range msg.Answer {
		if ls.MinTTL != 0 {
			if a.Header().Ttl < ls.MinTTL {
				a.Header().Ttl = ls.MinTTL
			}
		}

		if ls.MaxTTL != 0 {
			if a.Header().Ttl > ls.MaxTTL {
				a.Header().Ttl = ls.MaxTTL
			}
		}

	}
}

func (ls LocalServer) ListenAndServe() {
	ls.cache = gcache.New(ls.cacheCount).LFU().Build()

	log.Logger.Infof("listen server %s port:%s ", ls.ServerType.String(), ls.Port)
	srv := &dns.Server{Addr: ":" + ls.Port, Net: ls.ServerType.String()}
	srv.Handler = &ls
	if err := srv.ListenAndServe(); err != nil {
		log.Logger.Errorf("Failed to set %s listener %s\n", ls.ServerType.String(), err.Error())
	}
}

func (ls *LocalServer) MatchRemoteServer(r *dns.Msg) (constant.RemoteDnsServer, error) {
	rule, err := ls.RulesList.Match(r)
	if err != nil {
		if len(r.Question) > 0 {
			domain := r.Question[0].Name
			log.Logger.Debugf("domain:[%s] There are no matching rules, use the default server:%s", domain, ls.forward.String())
		}

		return nil, err
	} else {
		if len(r.Question) > 0 {
			domain := r.Question[0].Name
			log.Logger.Debugf("domain:[%s] Match Rule Type %s RemoteServer:%s", domain, rule.RuleType(), rule.RemoteServer().String())
		}

	}
	remoteServer := rule.RemoteServer()
	return remoteServer, nil
}

func NewLocalServer(port string, serverType LocalServerType, minTTL, maxTTL uint32, cacheCount int, rulesList *rules.RuleList, forward constant.RemoteDnsServer) LocalServer {
	if cacheCount == 0 {
		cacheCount = constant.DefaultCacheCount
	}
	return LocalServer{
		Port:       port,
		ServerType: serverType,
		MinTTL:     minTTL,
		MaxTTL:     maxTTL,
		RulesList:  rulesList,
		cacheCount: cacheCount,
		forward:    forward,
	}
}
